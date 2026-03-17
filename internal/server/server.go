package server

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/games"
	"github.com/jgoodhcg/mindmeld/internal/games/cluster"
	"github.com/jgoodhcg/mindmeld/internal/games/trivia"
	"github.com/jgoodhcg/mindmeld/internal/ws"
)

type Server struct {
	queries  *db.Queries
	dbPool   *pgxpool.Pool
	router   *chi.Mux
	hub      *ws.Hub
	eventBus events.Bus
	games    *games.Registry
}

func NewServer(pool *pgxpool.Pool) *Server {
	hub := ws.NewHub()
	eventBus := events.NewInMemoryBus()
	queries := db.New(pool)

	// Create game registry and register games
	registry := games.NewRegistry()
	triviaGame := trivia.New(queries, pool, eventBus, hub)
	clusterGame := cluster.New(queries, pool, eventBus, hub)
	registry.Register(triviaGame)
	registry.Register(clusterGame)

	// Register placeholder games (coming soon)
	registry.RegisterPlaceholder(games.GameInfo{Slug: "wavelength", Name: "WAVELENGTH", Description: "Find the spectrum between extremes"})
	registry.RegisterPlaceholder(games.GameInfo{Slug: "sync", Name: "SYNC", Description: "Answer simultaneously, match to score"})
	registry.RegisterPlaceholder(games.GameInfo{Slug: "cipher", Name: "CIPHER", Description: "One word clues, shared understanding"})

	s := &Server{
		queries:  queries,
		dbPool:   pool,
		router:   chi.NewRouter(),
		hub:      hub,
		eventBus: eventBus,
		games:    registry,
	}

	hub.SetPresenceHandler(func(lobbyCode string, update ws.PresenceUpdate) {
		eventBus.Publish(context.Background(), events.Event{
			Type:      events.EventPlayerPresence,
			LobbyCode: lobbyCode,
			Payload: events.PlayerPresencePayload{
				PlayerID:     update.PlayerID,
				Connected:    update.Connected,
				GraceExpired: update.GraceExpired,
			},
		})
	})

	// Wire up event subscriber to broadcast updates
	subscriber := ws.NewSubscriber(hub, queries, registry, pool)
	eventBus.Subscribe(subscriber.HandleEvent)

	s.routes()
	return s
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) SetDisconnectGracePeriod(gracePeriod time.Duration) {
	s.hub.SetDisconnectGracePeriod(gracePeriod)
}
