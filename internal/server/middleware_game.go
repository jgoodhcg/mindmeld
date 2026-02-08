package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// requireLobbyGameType enforces the contract for game-scoped endpoints:
// /lobbies/{code}/{game}/... must only operate on lobbies whose current
// game_type matches {game}. This keeps game actions isolated even if a lobby
// can switch game types over time.
func (s *Server) requireLobbyGameType(expectedGameType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := chi.URLParam(r, "code")
			lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
			if err != nil {
				http.Error(w, "Lobby not found", http.StatusNotFound)
				return
			}

			if !strings.EqualFold(lobby.GameType, expectedGameType) {
				log.Printf("Blocked game-scoped action: lobby=%s lobby_game_type=%s route_game_type=%s", code, lobby.GameType, expectedGameType)
				http.Error(w, "Lobby is currently set to a different game", http.StatusConflict)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
