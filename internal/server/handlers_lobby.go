package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/templates"
)

func (s *Server) handleCreateLobby(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	lobbyName := strings.TrimSpace(r.FormValue("name"))
	nickname := strings.TrimSpace(r.FormValue("nickname"))

	if lobbyName == "" || nickname == "" {
		http.Error(w, "Lobby name and nickname are required", http.StatusBadRequest)
		return
	}

	player := GetPlayer(r.Context())
	code := generateCode()

	// Create Lobby
	lobby, err := s.queries.CreateLobby(r.Context(), db.CreateLobbyParams{
		Code: code,
		Name: lobbyName,
	})
	if err != nil {
		log.Printf("Error creating lobby: %v", err)
		http.Error(w, "Failed to create lobby", http.StatusInternalServerError)
		return
	}

	// Add Host Player
	_, err = s.queries.AddPlayerToLobby(r.Context(), db.AddPlayerToLobbyParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
		Nickname: nickname,
		IsHost:   true,
	})
	if err != nil {
		log.Printf("Error adding host to lobby: %v", err)
		http.Error(w, "Failed to join lobby", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/lobbies/"+lobby.Code, http.StatusSeeOther)
}

func (s *Server) handleLobbyRoom(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	
	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	players, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		log.Printf("Error fetching lobby players: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templates.LobbyRoom(lobby, players).Render(r.Context(), w)
}
