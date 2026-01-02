package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

type Server struct {
	queries *db.Queries
	router  *chi.Mux
}

func NewServer(queries *db.Queries) *Server {
	s := &Server{
		queries: queries,
		router:  chi.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
