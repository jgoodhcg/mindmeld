package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/ws"
)

type Server struct {
	queries  *db.Queries
	router   *chi.Mux
	hub      *ws.Hub
	eventBus events.Bus
}

func NewServer(queries *db.Queries) *Server {
	hub := ws.NewHub()
	eventBus := events.NewInMemoryBus()

	s := &Server{
		queries:  queries,
		router:   chi.NewRouter(),
		hub:      hub,
		eventBus: eventBus,
	}

	// Wire up event subscriber to broadcast updates
	subscriber := ws.NewSubscriber(hub, queries)
	eventBus.Subscribe(subscriber.HandleEvent)

	s.routes()
	return s
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
