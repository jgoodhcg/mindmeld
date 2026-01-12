package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/jgoodhcg/mindmeld/templates"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	stats, err := s.queries.GetLobbyStats(r.Context())
	if err != nil {
		log.Printf("Error getting lobby stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	templates.Home(stats.TotalLobbies, stats.LobbiesWithPlayers).Render(r.Context(), w)
}

func (s *Server) handleJoinByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(strings.TrimSpace(r.FormValue("code")))
	if code == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}
