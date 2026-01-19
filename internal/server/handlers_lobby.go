package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	player := GetPlayer(r.Context())
	
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

	// If Playing, get active round
	var activeRound db.TriviaRound
	var hasSubmitted bool
	var currentQuestion db.TriviaQuestion
	var questionActive bool
	var isAuthor bool
	var hasAnswered bool
	var submittedCount int
	var scoreboard []db.GetLobbyScoreboardRow
	var distribution []events.AnswerStat
	var totalAnswers int
	
	if lobby.Phase == "playing" {
		activeRound, err = s.queries.GetActiveRound(r.Context(), lobby.ID)
		if err != nil {
			log.Printf("Error fetching active round: %v", err)
		} else {
			if activeRound.Phase == "submitting" {
				// Check if player submitted AND count total submissions
				questions, err := s.queries.GetQuestionsForRound(r.Context(), activeRound.ID)
				if err == nil {
					submittedCount = len(questions)
					for _, q := range questions {
						if q.Author == player.ID {
							hasSubmitted = true
						}
					}
				}
			} else if activeRound.Phase == "playing" {
				// Use Explicit State
				if activeRound.CurrentQuestionID.Valid {
					// Get the current question
					var qID pgtype.UUID = activeRound.CurrentQuestionID
					questions, err := s.queries.GetQuestionsForRound(r.Context(), activeRound.ID)
					if err == nil {
						for _, q := range questions {
							if q.ID == qID {
								currentQuestion = q
								questionActive = true
								break
							}
						}
					}

					if questionActive {
						// Check if Author
						if currentQuestion.Author == player.ID {
							isAuthor = true
						}

						// Check if Answered
						answers, err := s.queries.GetAnswersForQuestion(r.Context(), currentQuestion.ID)
						if err == nil {
							totalAnswers = len(answers)
							for _, a := range answers {
								if a.PlayerID == player.ID {
									hasAnswered = true
								}
								// Build Stats
								found := false
								for i, stat := range distribution {
									if stat.Answer == a.SelectedAnswer {
										distribution[i].Count++
										found = true
										break
									}
								}
								if !found {
									distribution = append(distribution, events.AnswerStat{
										Answer: a.SelectedAnswer,
										Count:  1,
									})
								}
							}
						}
					}
				} else {
					// No current question set? Should be finished or error.
					// Fallback to checking if round is actually finished
				}
			} else if activeRound.Phase == "finished" {
				scoreboard, err = s.queries.GetLobbyScoreboard(r.Context(), lobby.ID)
				if err != nil {
					log.Printf("Error fetching scoreboard: %v", err)
				}
			}
		}
	}

	templates.LobbyRoom(lobby, players, activeRound, hasSubmitted, currentQuestion, questionActive, isAuthor, hasAnswered, submittedCount, participation.IsHost, scoreboard, distribution, totalAnswers).Render(r.Context(), w)
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

	player := GetPlayer(r.Context())

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
