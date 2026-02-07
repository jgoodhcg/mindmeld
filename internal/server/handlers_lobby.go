package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/templates"
)

func (s *Server) handleCreateLobby(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	lobbyName := strings.TrimSpace(r.FormValue("name"))
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	gameType := strings.TrimSpace(r.FormValue("game_type"))

	if lobbyName == "" || nickname == "" {
		http.Error(w, "Lobby name and nickname are required", http.StatusBadRequest)
		return
	}

	if gameType == "" {
		gameType = "trivia"
	}

	player := auth.GetPlayer(r.Context())
	code := generateCode()

	// Create Lobby
	lobby, err := s.queries.CreateLobby(r.Context(), db.CreateLobbyParams{
		Code:     code,
		Name:     lobbyName,
		GameType: gameType,
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

	// Publish event for real-time updates (host joining counts as a player join)
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventPlayerJoined,
		LobbyCode: code,
		Payload: events.PlayerJoinedPayload{
			PlayerID: player.ID.String(),
			Nickname: nickname,
		},
	})

	http.Redirect(w, r, "/lobbies/"+lobby.Code, http.StatusSeeOther)
}

func (s *Server) handleLobbyRoom(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	// Check if player is already participating
	participation, err := s.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})

	// If checking failed (not found), show Join screen
	if err != nil {
		templates.JoinLobby(lobby).Render(r.Context(), w)
		return
	}

	// Player is in the lobby, show the room
	players, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		log.Printf("Error fetching lobby players: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Look up game by lobby's game_type and render content
	game, ok := s.games.Get(lobby.GameType)
	if !ok {
		log.Printf("Unknown game type: %s", lobby.GameType)
		http.Error(w, "Unknown game type", http.StatusInternalServerError)
		return
	}

	gameContent := game.RenderContent(r.Context(), lobby, players, player, participation.IsHost)
	templates.LobbyRoom(lobby, players, participation.IsHost, gameContent).Render(r.Context(), w)
}

func (s *Server) handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	nickname := strings.TrimSpace(r.FormValue("nickname"))

	if nickname == "" {
		http.Error(w, "Nickname is required", http.StatusBadRequest)
		return
	}

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Add player to lobby
	_, err = s.queries.AddPlayerToLobby(r.Context(), db.AddPlayerToLobbyParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
		Nickname: nickname,
		IsHost:   false,
	})
	if err != nil {
		// If duplicate key error, maybe they refreshed or clicked twice?
		// For now, log and fail, or we could redirect if they are already in.
		log.Printf("Error joining lobby: %v", err)
		http.Error(w, "Failed to join lobby (Name taken?)", http.StatusInternalServerError)
		return
	}

	// Publish event for real-time updates
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventPlayerJoined,
		LobbyCode: code,
		Payload: events.PlayerJoinedPayload{
			PlayerID: player.ID.String(),
			Nickname: nickname,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}
