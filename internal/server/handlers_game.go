package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

func (s *Server) handleGetGameContent(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	participation, err := s.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil {
		http.Error(w, "Not in lobby", http.StatusForbidden)
		return
	}

	players, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	game, ok := s.games.Get(lobby.GameType)
	if !ok {
		log.Printf("Unknown game type: %s", lobby.GameType)
		http.Error(w, "Unknown game type", http.StatusInternalServerError)
		return
	}

	game.RenderContent(r.Context(), lobby, players, player, participation.IsHost).Render(r.Context(), w)
}
