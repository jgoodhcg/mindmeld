package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

func (s *Server) handleStartGame(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := GetPlayer(r.Context())

	// Verify Host
	participation, err := s.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil {
		log.Printf("Error checking host status: %v", err)
		http.Error(w, "Error checking permissions", http.StatusInternalServerError)
		return
	}
	if !participation.IsHost {
		log.Printf("Player %v is not host of lobby %v", player.ID, lobby.ID)
		http.Error(w, "Only the host can start the game", http.StatusForbidden)
		return
	}

	log.Printf("Starting game for lobby %s", lobby.Code)

	// 1. Update Lobby Phase
	err = s.queries.UpdateLobbyPhase(r.Context(), db.UpdateLobbyPhaseParams{
		ID:    lobby.ID,
		Phase: "playing",
	})
	if err != nil {
		log.Printf("Error updating lobby phase: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	
	log.Printf("Lobby phase updated to playing")

	// 2. Create Round 1
	// We use the default phase 'submitting' defined in SQL schema
	_, err = s.queries.CreateTriviaRound(r.Context(), db.CreateTriviaRoundParams{
		LobbyID:     lobby.ID,
		RoundNumber: 1,
	})
	if err != nil {
		log.Printf("Error creating round: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (s *Server) handleSubmitQuestion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := GetPlayer(r.Context())
	
	// Get active round
	round, err := s.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	// Create Question
	_, err = s.queries.CreateQuestion(r.Context(), db.CreateQuestionParams{
		RoundID:       round.ID,
		Author:        player.ID,
		QuestionText:  r.FormValue("question_text"),
		CorrectAnswer: r.FormValue("correct_answer"),
		WrongAnswer1:  r.FormValue("wrong_answer_1"),
		WrongAnswer2:  r.FormValue("wrong_answer_2"),
		WrongAnswer3:  r.FormValue("wrong_answer_3"),
	})
	if err != nil {
		log.Printf("Error creating question: %v", err)
		http.Error(w, "Failed to submit question", http.StatusInternalServerError)
		return
	}

	// TODO: Check if all players have submitted to auto-advance?
	// For now, just redirect back to lobby
	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}
