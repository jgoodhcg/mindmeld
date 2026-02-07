package ws

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
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
}

// NewSubscriber creates a new event subscriber.
func NewSubscriber(hub *Hub, queries *db.Queries, registry GameRegistry) *Subscriber {
	return &Subscriber{
		hub:      hub,
		queries:  queries,
		registry: registry,
	}
}

// HandleEvent processes an event and broadcasts appropriate updates.
func (s *Subscriber) HandleEvent(ctx context.Context, event events.Event) {
	switch event.Type {
	case events.EventPlayerJoined, events.EventPlayerLeft:
		s.broadcastPlayerList(ctx, event.LobbyCode)
	default:
		// Look up the lobby's game type and delegate to the game handler
		lobby, err := s.queries.GetLobbyByCode(ctx, event.LobbyCode)
		if err != nil {
			log.Printf("[ws-subscriber] Failed to get lobby %s: %v", event.LobbyCode, err)
			return
		}

		handler, ok := s.registry.GetHandler(lobby.GameType)
		if !ok {
			log.Printf("[ws-subscriber] No handler for game type: %s", lobby.GameType)
			return
		}

		if !handler.HandleEvent(ctx, event, s.hub, s.queries) {
			log.Printf("[ws-subscriber] Unhandled event type: %s", event.Type)
		}
	}
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
	err = templates.PlayerList(players, true).Render(ctx, &buf)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to render player list: %v", err)
		return
	}

	s.hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}
