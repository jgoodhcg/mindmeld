package ws

import (
	"bytes"
	"context"
	"log"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/templates"
)

// Subscriber listens to events and broadcasts updates to WebSocket clients.
type Subscriber struct {
	hub     *Hub
	queries *db.Queries
}

// NewSubscriber creates a new event subscriber.
func NewSubscriber(hub *Hub, queries *db.Queries) *Subscriber {
	return &Subscriber{
		hub:     hub,
		queries: queries,
	}
}

// HandleEvent processes an event and broadcasts appropriate updates.
func (s *Subscriber) HandleEvent(ctx context.Context, event events.Event) {
	switch event.Type {
	case events.EventPlayerJoined, events.EventPlayerLeft:
		s.broadcastPlayerList(ctx, event.LobbyCode)
	// Future event types can be handled here:
	// case events.EventGameStarted:
	// case events.EventQuestionSubmitted:
	// case events.EventAnswerSubmitted:
	default:
		log.Printf("[ws-subscriber] Unhandled event type: %s", event.Type)
	}
}

// broadcastPlayerList fetches the current player list and broadcasts it.
func (s *Subscriber) broadcastPlayerList(ctx context.Context, lobbyCode string) {
	// Get the lobby to find its ID
	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get lobby %s: %v", lobbyCode, err)
		return
	}

	// Fetch current players
	players, err := s.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get players for lobby %s: %v", lobbyCode, err)
		return
	}

	// Render the player list partial to HTML
	var buf bytes.Buffer
	err = templates.PlayerList(players).Render(ctx, &buf)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to render player list: %v", err)
		return
	}

	// Broadcast to all connected clients in this lobby
	s.hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}
