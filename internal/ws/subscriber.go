package ws

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/lobbyview"
	"github.com/jgoodhcg/mindmeld/templates"
)

// GameHandler is the interface that game types implement for event handling.
// This avoids an import cycle with the games package.
type GameHandler interface {
	HandleEvent(ctx context.Context, event events.Event, hub *Hub, queries *db.Queries) bool
}

// GameRegistry provides game lookup by slug.
type GameRegistry interface {
	GetHandler(gameType string) (GameHandler, bool)
}

// Subscriber listens to events and broadcasts updates to WebSocket clients.
type Subscriber struct {
	hub      *Hub
	queries  *db.Queries
	registry GameRegistry
	dbPool   *pgxpool.Pool
}

// NewSubscriber creates a new event subscriber.
func NewSubscriber(hub *Hub, queries *db.Queries, registry GameRegistry, dbPool *pgxpool.Pool) *Subscriber {
	return &Subscriber{
		hub:      hub,
		queries:  queries,
		registry: registry,
		dbPool:   dbPool,
	}
}

// HandleEvent processes an event and broadcasts appropriate updates.
func (s *Subscriber) HandleEvent(ctx context.Context, event events.Event) {
	switch event.Type {
	case events.EventPlayerJoined, events.EventPlayerLeft, events.EventHostTransferred:
		s.broadcastPlayerList(ctx, event.LobbyCode)
		// Player changes can affect game-content state (e.g. minimum players/start button).
		// Trigger a full content refresh in addition to the player-list OOB swap.
		BroadcastUpdateTrigger(ctx, event.LobbyCode, s.hub)
	case events.EventPlayerPresence:
		payload, ok := event.Payload.(events.PlayerPresencePayload)
		if ok && payload.GraceExpired {
			s.maybeTransferHost(ctx, event.LobbyCode)
		}
		s.broadcastPlayerList(ctx, event.LobbyCode)
		s.delegateToGame(ctx, event)
		if !ok || s.shouldRefreshGameContentForPresence(ctx, event.LobbyCode, payload) {
			BroadcastUpdateTrigger(ctx, event.LobbyCode, s.hub)
		}
	default:
		if !s.delegateToGame(ctx, event) {
			log.Printf("[ws-subscriber] Unhandled event type: %s", event.Type)
		}
	}
}

func (s *Subscriber) delegateToGame(ctx context.Context, event events.Event) bool {
	lobby, err := s.queries.GetLobbyByCode(ctx, event.LobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get lobby %s: %v", event.LobbyCode, err)
		return false
	}

	handler, ok := s.registry.GetHandler(lobby.GameType)
	if !ok {
		log.Printf("[ws-subscriber] No handler for game type: %s", lobby.GameType)
		return false
	}

	return handler.HandleEvent(ctx, event, s.hub, s.queries)
}

func (s *Subscriber) shouldRefreshGameContentForPresence(ctx context.Context, lobbyCode string, payload events.PlayerPresencePayload) bool {
	if payload.GraceExpired {
		return true
	}

	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to load lobby %s for presence refresh check: %v", lobbyCode, err)
		return true
	}

	// During active Cluster rounds, transient reconnect churn from another player's form POST
	// should not reset in-progress coordinate selection for everyone else.
	if lobby.GameType == "cluster" && strings.EqualFold(lobby.Phase, "playing") {
		return false
	}

	return true
}

// BroadcastUpdateTrigger sends an OOB swap that triggers the client to fetch updated game content.
// Exported for use by game packages.
func BroadcastUpdateTrigger(ctx context.Context, lobbyCode string, hub *Hub) {
	html := fmt.Sprintf(`<div id="game-updater" hx-swap-oob="true" hx-get="/lobbies/%s/content" hx-target="#game-content" hx-swap="outerHTML" hx-trigger="load" class="hidden"></div>`, lobbyCode)
	hub.Broadcast(ctx, lobbyCode, []byte(html))
}

// broadcastPlayerList fetches the current player list and broadcasts it.
func (s *Subscriber) broadcastPlayerList(ctx context.Context, lobbyCode string) {
	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get lobby %s: %v", lobbyCode, err)
		return
	}

	players, err := s.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get players for lobby %s: %v", lobbyCode, err)
		return
	}

	var buf bytes.Buffer
	now := time.Now()
	rawPresence := s.hub.Snapshot(lobbyCode)
	presence := make(map[string]lobbyview.Presence, len(rawPresence))
	for playerID, state := range rawPresence {
		presence[playerID] = lobbyview.Presence{
			Disconnected: state.IsDisconnected(),
			GraceExpired: state.GraceExpiredAt(now),
		}
	}
	playerViews := lobbyview.Build(players, presence)
	err = templates.PlayerList(playerViews, true).Render(ctx, &buf)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to render player list: %v", err)
		return
	}

	s.hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}

func (s *Subscriber) maybeTransferHost(ctx context.Context, lobbyCode string) {
	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to load lobby %s for host transfer: %v", lobbyCode, err)
		return
	}

	players, err := s.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to load players for host transfer in lobby %s: %v", lobbyCode, err)
		return
	}

	now := time.Now()
	var currentHost *db.GetLobbyPlayersRow
	for i := range players {
		if players[i].IsHost {
			currentHost = &players[i]
			break
		}
	}
	if currentHost == nil {
		return
	}
	if s.hub.Presence(lobbyCode, currentHost.PlayerID.String()).IsActiveAt(now) {
		return
	}

	var nextHost *db.GetLobbyPlayersRow
	for i := range players {
		if players[i].PlayerID == currentHost.PlayerID {
			continue
		}
		if s.hub.Presence(lobbyCode, players[i].PlayerID.String()).IsConnected() {
			nextHost = &players[i]
			break
		}
	}
	if nextHost == nil {
		return
	}

	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to begin host transfer for lobby %s: %v", lobbyCode, err)
		return
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `UPDATE lobby_players SET is_host = false WHERE lobby_id = $1`, lobby.ID); err != nil {
		log.Printf("[ws-subscriber] Failed clearing host in lobby %s: %v", lobbyCode, err)
		return
	}
	if _, err = tx.Exec(ctx, `UPDATE lobby_players SET is_host = true WHERE lobby_id = $1 AND player_id = $2`, lobby.ID, nextHost.PlayerID); err != nil {
		log.Printf("[ws-subscriber] Failed assigning host in lobby %s: %v", lobbyCode, err)
		return
	}
	if err = tx.Commit(ctx); err != nil {
		log.Printf("[ws-subscriber] Failed committing host transfer in lobby %s: %v", lobbyCode, err)
		return
	}

	log.Printf("[ws-subscriber] Transferred host in lobby %s from %s to %s", lobbyCode, currentHost.Nickname, nextHost.Nickname)
}
