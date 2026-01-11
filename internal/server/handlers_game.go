package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
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

	// 3. Publish game started event for real-time updates
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventGameStarted,
		LobbyCode: code,
		Payload: events.GameStartedPayload{
			RoundNumber: 1,
		},
	})

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

	// Get updated counts for real-time update
	players, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		log.Printf("Error getting players: %v", err)
	}
	questions, err := s.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err != nil {
		log.Printf("Error getting questions: %v", err)
	}

	// Find the host player ID
	var hostPlayerID string
	for _, p := range players {
		if p.IsHost {
			hostPlayerID = p.PlayerID.String()
			break
		}
	}

	// Publish question submitted event for real-time updates
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventQuestionSubmitted,
		LobbyCode: code,
		Payload: events.QuestionSubmittedPayload{
			SubmittedCount: len(questions),
			TotalPlayers:   len(players),
			HostPlayerID:   hostPlayerID,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (s *Server) handleAdvanceRound(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can advance the round", http.StatusForbidden)
		return
	}

	// Get active round
	round, err := s.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	// Logic moved from handleSubmitQuestion:
	questions, err := s.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err == nil {
		// 1. Shuffle and assign order
		for i, q := range questions {
			err := s.queries.UpdateQuestionOrder(r.Context(), db.UpdateQuestionOrderParams{
				ID:           q.ID,
				DisplayOrder: pgtype.Int4{Int32: int32(i + 1), Valid: true},
			})
			if err != nil {
				log.Printf("Error updating order for q %v: %v", q.ID, err)
			}
		}

		// 2. Advance Round Phase
		err = s.queries.UpdateRoundPhase(r.Context(), db.UpdateRoundPhaseParams{
			ID:    round.ID,
			Phase: "playing",
		})
		if err != nil {
			log.Printf("Error advancing round phase: %v", err)
		}

		// Publish round advanced event for real-time updates
		s.eventBus.Publish(r.Context(), events.Event{
			Type:      events.EventRoundAdvanced,
			LobbyCode: code,
			Payload: events.RoundAdvancedPayload{
				RoundNumber: round.RoundNumber,
			},
		})
	}

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (s *Server) handlePlayAgain(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can start a new round", http.StatusForbidden)
		return
	}

	// Get current active round to increment number
	lastRound, err := s.queries.GetActiveRound(r.Context(), lobby.ID)
	nextRoundNum := int32(1)
	if err == nil {
		nextRoundNum = lastRound.RoundNumber + 1
	}

	// Create New Round
	_, err = s.queries.CreateTriviaRound(r.Context(), db.CreateTriviaRoundParams{
		LobbyID:     lobby.ID,
		RoundNumber: nextRoundNum,
	})
	if err != nil {
		log.Printf("Error creating new round: %v", err)
		http.Error(w, "Failed to start new round", http.StatusInternalServerError)
		return
	}

	// Publish new round event for real-time updates
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventNewRoundCreated,
		LobbyCode: code,
		Payload: events.NewRoundCreatedPayload{
			RoundNumber: nextRoundNum,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (s *Server) handleSubmitAnswer(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	questionIDStr := r.FormValue("question_id")
	selectedAnswer := r.FormValue("answer")

	// Parse UUID
	var questionID pgtype.UUID
	if err := questionID.Scan(questionIDStr); err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// 1. Get Lobby
	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}
	
	// 2. Get Active Round
	round, err := s.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	// 3. Get Questions to verify it exists in this round
	questions, err := s.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err != nil {
		http.Error(w, "Error fetching questions", http.StatusInternalServerError)
		return
	}

	var targetQuestion db.TriviaQuestion
	found := false
	for _, q := range questions {
		if q.ID == questionID {
			targetQuestion = q
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Question not found in active round", http.StatusBadRequest)
		return
	}

	// 4. Check Answer
	isCorrect := (selectedAnswer == targetQuestion.CorrectAnswer)
	player := GetPlayer(r.Context())

	// 5. Submit Answer
	_, err = s.queries.SubmitAnswer(r.Context(), db.SubmitAnswerParams{
		QuestionID:     questionID,
		PlayerID:       player.ID,
		SelectedAnswer: selectedAnswer,
		IsCorrect:      isCorrect,
	})
	if err != nil {
		log.Printf("Error submitting answer: %v", err)
		// Ignore error (maybe duplicate submission)
	}

	// Check if round is finished (all questions answered by all other players)
	// We reuse logic similar to lobby room check, or just check global counts.
	// Simple check: do ANY questions remain with < (players-1) answers?

	lobbyPlayers, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	roundFinished := false
	totalAnswered := 0
	totalExpected := 0

	if err == nil {
		allFinished := true
		questions, err := s.queries.GetQuestionsForRound(r.Context(), round.ID)
		if err == nil {
			// Calculate totals for the event payload
			numPlayers := len(lobbyPlayers)
			numQuestions := len(questions)
			// Each question can be answered by (players - 1) since author doesn't answer their own
			totalExpected = numQuestions * (numPlayers - 1)
			if totalExpected < 0 {
				totalExpected = 0
			}

			for _, q := range questions {
				count, err := s.queries.CountAnswersForQuestion(r.Context(), q.ID)
				if err != nil {
					continue
				}
				totalAnswered += int(count)

				target := int64(numPlayers - 1)
				if target < 0 {
					target = 0
				}

				if count < target {
					allFinished = false
				}
			}
		}

		if allFinished {
			roundFinished = true
			err = s.queries.UpdateRoundPhase(r.Context(), db.UpdateRoundPhaseParams{
				ID:    round.ID,
				Phase: "finished",
			})
			if err != nil {
				log.Printf("Error finishing round: %v", err)
			}
		}
	}

	// Publish answer submitted event for real-time updates
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventAnswerSubmitted,
		LobbyCode: code,
		Payload: events.AnswerSubmittedPayload{
			AnsweredCount: totalAnswered,
			TotalExpected: totalExpected,
			RoundFinished: roundFinished,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}