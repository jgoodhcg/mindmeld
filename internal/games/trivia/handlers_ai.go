package trivia

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

func (g *TriviaGame) handleGenerateQuestion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	_, err = g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil {
		http.Error(w, "Not in lobby", http.StatusForbidden)
		return
	}

	round, err := g.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil || !strings.EqualFold(round.Phase, "submitting") {
		http.Error(w, "Question generation is only available during submission", http.StatusConflict)
		return
	}

	used, err := g.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err == nil {
		for _, q := range used {
			if q.Author == player.ID {
				http.Error(w, "You already submitted this round", http.StatusConflict)
				return
			}
		}
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	topic := strings.TrimSpace(r.FormValue("topic"))
	generated, genErr := generateAssistedQuestion(r.Context(), lobby.ContentRating, topic)
	if genErr != nil {
		log.Printf("Error generating question for lobby %s: %v", lobby.Code, genErr)
		writeGenerateQuestionResponse(w, http.StatusInternalServerError, generateQuestionResponse{
			Error: "Unable to generate a question right now.",
		})
		return
	}

	writeGenerateQuestionResponse(w, http.StatusOK, generateQuestionResponse{
		QuestionText:  generated.QuestionText,
		CorrectAnswer: generated.CorrectAnswer,
		WrongAnswer1:  generated.WrongAnswer1,
		WrongAnswer2:  generated.WrongAnswer2,
		WrongAnswer3:  generated.WrongAnswer3,
		Source:        generated.Source,
	})
}

func writeGenerateQuestionResponse(w http.ResponseWriter, status int, payload generateQuestionResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
