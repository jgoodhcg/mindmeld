package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/jgoodhcg/mindmeld/internal/db"
)

func (s *Server) handleCreateLobby(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	code := generateCode()

	_, err := s.queries.CreateLobby(r.Context(), db.CreateLobbyParams{
		Code: code,
		Name: name,
	})
	if err != nil {
		log.Printf("Error creating lobby: %v", err)
		http.Error(w, "Failed to create lobby", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
