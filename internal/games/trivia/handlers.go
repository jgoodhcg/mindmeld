package trivia

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/questions"
	triviatmpl "github.com/jgoodhcg/mindmeld/templates/trivia"
)

func (g *TriviaGame) handleStartGame(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Verify Host
	participation, err := g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
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
	err = g.queries.UpdateLobbyPhase(r.Context(), db.UpdateLobbyPhaseParams{
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
	_, err = g.queries.CreateTriviaRound(r.Context(), db.CreateTriviaRoundParams{
		LobbyID:     lobby.ID,
		RoundNumber: 1,
	})
	if err != nil {
		log.Printf("Error creating round: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// 3. Publish game started event for real-time updates
	g.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventGameStarted,
		LobbyCode: code,
		Payload: events.GameStartedPayload{
			RoundNumber: 1,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *TriviaGame) handleGetQuestionTemplates(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	// Verify player is in lobby
	_, err = g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil {
		http.Error(w, "Not in lobby", http.StatusForbidden)
		return
	}

	// Get used templates for this lobby
	usedTemplateIDs, err := g.queries.GetUsedTemplatesForLobby(r.Context(), lobby.ID)
	if err != nil {
		log.Printf("Error getting used templates: %v", err)
		usedTemplateIDs = []string{}
	}

	// Get available templates (not yet used)
	available := questions.GetAvailableTemplates(usedTemplateIDs)
	grouped := questions.GroupByCategory(available)
	categories := questions.GetCategories()

	triviatmpl.QuestionTemplatesModal(grouped, categories).Render(r.Context(), w)
}

func (g *TriviaGame) handleSubmitQuestion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Get active round
	round, err := g.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tx, err := g.dbPool.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Failed to submit question", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	qtx := g.queries.WithTx(tx)

	// Create Question
	_, err = qtx.CreateQuestion(ctx, db.CreateQuestionParams{
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

	// Check if a template was used and mark it
	templateID := r.FormValue("template_id")
	if templateID != "" {
		err = qtx.MarkTemplateUsed(ctx, db.MarkTemplateUsedParams{
			LobbyID:    lobby.ID,
			TemplateID: templateID,
		})
		if err != nil {
			log.Printf("Error marking template as used: %v", err)
			http.Error(w, "Failed to submit question", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Failed to submit question", http.StatusInternalServerError)
		return
	}

	// Get updated counts for real-time update
	players, err := g.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		log.Printf("Error getting players: %v", err)
	}
	roundQuestions, err := g.queries.GetQuestionsForRound(r.Context(), round.ID)
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
	g.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventQuestionSubmitted,
		LobbyCode: code,
		Payload: events.QuestionSubmittedPayload{
			SubmittedCount: len(roundQuestions),
			TotalPlayers:   len(players),
			HostPlayerID:   hostPlayerID,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *TriviaGame) handleAdvanceRound(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Verify Host
	participation, err := g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can advance the round", http.StatusForbidden)
		return
	}

	// Get active round
	round, err := g.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	roundQuestions, err := g.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err == nil && len(roundQuestions) > 0 {
		// 1. Shuffle and assign order
		for i, q := range roundQuestions {
			err := g.queries.UpdateQuestionOrder(r.Context(), db.UpdateQuestionOrderParams{
				ID:           q.ID,
				DisplayOrder: pgtype.Int4{Int32: int32(i + 1), Valid: true},
			})
			if err != nil {
				log.Printf("Error updating order for q %v: %v", q.ID, err)
			}
		}

		// 2. Set Initial Question State
		firstQuestion := roundQuestions[0]
		err = g.queries.UpdateRoundQuestionState(r.Context(), db.UpdateRoundQuestionStateParams{
			ID:                round.ID,
			CurrentQuestionID: firstQuestion.ID,
			QuestionState:     "answering",
		})
		if err != nil {
			log.Printf("Error setting initial question state: %v", err)
		}

		// 3. Advance Round Phase
		err = g.queries.UpdateRoundPhase(r.Context(), db.UpdateRoundPhaseParams{
			ID:    round.ID,
			Phase: "playing",
		})
		if err != nil {
			log.Printf("Error advancing round phase: %v", err)
		}

		// Publish round advanced event for real-time updates
		g.eventBus.Publish(r.Context(), events.Event{
			Type:      events.EventRoundAdvanced,
			LobbyCode: code,
			Payload: events.RoundAdvancedPayload{
				RoundNumber: round.RoundNumber,
			},
		})
	}

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *TriviaGame) handleNextQuestion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Verify Host
	participation, err := g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can advance the question", http.StatusForbidden)
		return
	}

	// Get active round
	round, err := g.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	// Get Questions to find current and next
	roundQuestions, err := g.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err != nil {
		http.Error(w, "Error fetching questions", http.StatusInternalServerError)
		return
	}

	var nextQuestion *db.TriviaQuestion
	foundCurrent := false

	if !round.CurrentQuestionID.Valid {
		if len(roundQuestions) > 0 {
			nextQuestion = &roundQuestions[0]
		}
	} else {
		for i, q := range roundQuestions {
			if q.ID == round.CurrentQuestionID {
				foundCurrent = true
				if i+1 < len(roundQuestions) {
					nextQuestion = &roundQuestions[i+1]
				}
				break
			}
		}
	}

	if nextQuestion != nil {
		err = g.queries.UpdateRoundQuestionState(r.Context(), db.UpdateRoundQuestionStateParams{
			ID:                round.ID,
			CurrentQuestionID: nextQuestion.ID,
			QuestionState:     "answering",
		})
		if err != nil {
			log.Printf("Error updating next question: %v", err)
		}

		g.eventBus.Publish(r.Context(), events.Event{
			Type:      events.EventRoundAdvanced,
			LobbyCode: code,
			Payload: events.RoundAdvancedPayload{
				RoundNumber: round.RoundNumber,
			},
		})
	} else if foundCurrent {
		err = g.queries.UpdateRoundPhase(r.Context(), db.UpdateRoundPhaseParams{
			ID:    round.ID,
			Phase: "finished",
		})
		if err != nil {
			log.Printf("Error finishing round: %v", err)
		}

		g.eventBus.Publish(r.Context(), events.Event{
			Type:      events.EventRoundAdvanced,
			LobbyCode: code,
			Payload: events.RoundAdvancedPayload{
				RoundNumber: round.RoundNumber,
			},
		})
	}

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *TriviaGame) handlePlayAgain(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Verify Host
	participation, err := g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can start a new round", http.StatusForbidden)
		return
	}

	// Get current active round to increment number
	lastRound, err := g.queries.GetActiveRound(r.Context(), lobby.ID)
	nextRoundNum := int32(1)
	if err == nil {
		nextRoundNum = lastRound.RoundNumber + 1
	}

	// Create New Round
	_, err = g.queries.CreateTriviaRound(r.Context(), db.CreateTriviaRoundParams{
		LobbyID:     lobby.ID,
		RoundNumber: nextRoundNum,
	})
	if err != nil {
		log.Printf("Error creating new round: %v", err)
		http.Error(w, "Failed to start new round", http.StatusInternalServerError)
		return
	}

	// Publish new round event for real-time updates
	g.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventNewRoundCreated,
		LobbyCode: code,
		Payload: events.NewRoundCreatedPayload{
			RoundNumber: nextRoundNum,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *TriviaGame) handleSubmitAnswer(w http.ResponseWriter, r *http.Request) {
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
	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	// 2. Get Active Round
	round, err := g.queries.GetActiveRound(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "No active round", http.StatusBadRequest)
		return
	}

	// 3. Get Questions to verify it exists in this round
	roundQuestions, err := g.queries.GetQuestionsForRound(r.Context(), round.ID)
	if err != nil {
		http.Error(w, "Error fetching questions", http.StatusInternalServerError)
		return
	}

	var targetQuestion db.TriviaQuestion
	found := false
	for _, q := range roundQuestions {
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
	player := auth.GetPlayer(r.Context())

	// 5. Submit Answer
	_, err = g.queries.SubmitAnswer(r.Context(), db.SubmitAnswerParams{
		QuestionID:     questionID,
		PlayerID:       player.ID,
		SelectedAnswer: selectedAnswer,
		IsCorrect:      isCorrect,
	})
	if err != nil {
		log.Printf("Error submitting answer: %v", err)
	}

	// 6. Calculate Stats & Check Completion
	lobbyPlayers, err := g.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	questionComplete := false

	// Prepare distribution stats
	distribution := []events.AnswerStat{}
	rawStats, err := g.queries.GetAnswerStats(r.Context(), questionID)
	if err == nil {
		for _, s := range rawStats {
			distribution = append(distribution, events.AnswerStat{
				Answer: s.SelectedAnswer,
				Count:  int(s.Count),
			})
		}
	}

	if err == nil {
		numPlayers := len(lobbyPlayers)
		targetPerQuestion := int64(numPlayers - 1) // Minus author
		if targetPerQuestion < 0 {
			targetPerQuestion = 0
		}

		currentCount, err := g.queries.CountAnswersForQuestion(r.Context(), questionID)
		if err == nil {
			if currentCount >= targetPerQuestion {
				questionComplete = true
			}
		}
	}

	// 7. Update State & Publish Events
	if questionComplete {
		err = g.queries.UpdateRoundQuestionState(r.Context(), db.UpdateRoundQuestionStateParams{
			ID:                round.ID,
			CurrentQuestionID: questionID,
			QuestionState:     "revealed",
		})
		if err != nil {
			log.Printf("Error updating question state to revealed: %v", err)
		}

		g.eventBus.Publish(r.Context(), events.Event{
			Type:      events.EventQuestionRevealed,
			LobbyCode: code,
			Payload: events.QuestionRevealedPayload{
				QuestionID:   questionIDStr,
				Distribution: distribution,
			},
		})
	}

	// Always publish update for live graphs (even if not complete)
	answeredCount := 0
	totalExpected := 0
	if len(lobbyPlayers) > 0 {
		totalExpected = len(lobbyPlayers) - 1 // minus author
		if cnt, err := g.queries.CountAnswersForQuestion(r.Context(), questionID); err == nil {
			answeredCount = int(cnt)
		}
	}
	g.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventAnswerSubmitted,
		LobbyCode: code,
		Payload: events.AnswerSubmittedPayload{
			AnsweredCount:    answeredCount,
			TotalExpected:    totalExpected,
			QuestionComplete: questionComplete,
			RoundFinished:    false,
			Distribution:     distribution,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}
