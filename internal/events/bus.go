// Package events provides an in-memory event bus for real-time updates.
//
// FUTURE SCALING: When horizontally scaling to multiple server instances,
// replace InMemoryBus with a PostgresBus implementation using LISTEN/NOTIFY.
// The Bus interface allows swapping implementations without changing game code.
//
// PostgresBus pattern:
//   - On startup: conn.Exec("LISTEN mindmeld_events")
//   - On publish: conn.Exec("NOTIFY mindmeld_events, $1", jsonPayload)
//   - Each server receives notifications and broadcasts to local WebSocket connections
package events

import (
	"context"
	"sync"
)

// Event represents something that happened in the system.
type Event struct {
	Type      string // e.g., "player.joined", "game.started", "question.submitted"
	LobbyCode string // The lobby this event relates to
	Payload   any    // Event-specific data
}

// EventHandler processes events.
type EventHandler func(ctx context.Context, event Event)

// Bus is the interface for publishing and subscribing to events.
// Implementations can be in-memory (single instance) or distributed (Postgres, Redis, etc.).
type Bus interface {
	// Publish sends an event to all subscribers.
	Publish(ctx context.Context, event Event)

	// Subscribe registers a handler to receive all events.
	Subscribe(handler EventHandler)
}

// InMemoryBus is a simple in-memory implementation of Bus.
// Suitable for single-instance deployments.
type InMemoryBus struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

// NewInMemoryBus creates a new in-memory event bus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make([]EventHandler, 0),
	}
}

// Publish sends an event to all registered handlers.
func (b *InMemoryBus) Publish(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers := make([]EventHandler, len(b.handlers))
	copy(handlers, b.handlers)
	b.mu.RUnlock()

	for _, handler := range handlers {
		// Run handlers synchronously for now.
		// Could be made async with goroutines if needed.
		handler(ctx, event)
	}
}

// Subscribe registers a handler to receive events.
func (b *InMemoryBus) Subscribe(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, handler)
}

// Event type constants
const (
	EventPlayerJoined      = "player.joined"
	EventPlayerLeft        = "player.left"
	EventGameStarted       = "game.started"
	EventQuestionSubmitted = "question.submitted"
	EventAnswerSubmitted   = "answer.submitted"
	EventQuestionRevealed  = "question.revealed" // Everyone answered, show correct answer
	EventRoundAdvanced     = "round.advanced"
	EventNewRoundCreated   = "round.created" // For "Play Again" - new round started
)

// PlayerJoinedPayload is the payload for EventPlayerJoined.
type PlayerJoinedPayload struct {
	PlayerID string
	Nickname string
}

// GameStartedPayload is the payload for EventGameStarted.
type GameStartedPayload struct {
	RoundNumber int32
}

// QuestionSubmittedPayload is the payload for EventQuestionSubmitted.
type QuestionSubmittedPayload struct {
	SubmittedCount int
	TotalPlayers   int
	HostPlayerID   string // The player ID (UUID string) of the lobby host
}

// RoundAdvancedPayload is the payload for EventRoundAdvanced.
type RoundAdvancedPayload struct {
	RoundNumber int32
}

// AnswerStats holds the count for a specific answer option.
type AnswerStat struct {
	Answer string
	Count  int
}

// AnswerSubmittedPayload is the payload for EventAnswerSubmitted.
type AnswerSubmittedPayload struct {
	AnsweredCount    int
	TotalExpected    int
	QuestionComplete bool // True when the current question has all expected answers
	RoundFinished    bool
	Distribution     []AnswerStat // Live stats for the graph
}

// QuestionRevealedPayload is the payload for EventQuestionRevealed.
type QuestionRevealedPayload struct {
	QuestionID   string
	Distribution []AnswerStat
}

// NewRoundCreatedPayload is the payload for EventNewRoundCreated (Play Again).
type NewRoundCreatedPayload struct {
	RoundNumber int32
}
