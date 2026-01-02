package server

import (
	"log"
	"net/http"

	"github.com/jgoodhcg/mindmeld/templates"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	lobbies, err := s.queries.ListLobbies(r.Context())
	if err != nil {
		log.Printf("Error listing lobbies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	templates.Home(lobbies).Render(r.Context(), w)
}
