package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/jgoodhcg/mindmeld/templates"
)

func (s *Server) handlePlatform(w http.ResponseWriter, r *http.Request) {
	allInfo := s.games.AllInfo()
	gameInfos := make([]templates.PlatformGameInfo, len(allInfo))
	for i, g := range allInfo {
		gameInfos[i] = templates.PlatformGameInfo{
			Slug:        g.Slug,
			Name:        g.Name,
			Description: g.Description,
			Ready:       g.Ready,
		}
	}
	templates.Platform(gameInfos).Render(r.Context(), w)
}

func (s *Server) handleTriviaHome(w http.ResponseWriter, r *http.Request) {
	stats, err := s.queries.GetLobbyStats(r.Context())
	if err != nil {
		log.Printf("Error getting lobby stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	liveLobbies := int64(s.hub.LiveLobbyCount())
	templates.Home(stats.TotalLobbies, stats.TotalAnswers, stats.TotalRounds, liveLobbies).Render(r.Context(), w)
}

func (s *Server) handleJoinByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(strings.TrimSpace(r.FormValue("code")))
	if code == "" {
		http.Redirect(w, r, "/trivia", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}
