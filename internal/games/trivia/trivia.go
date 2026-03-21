package trivia

import (
	"context"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/games"
	"github.com/jgoodhcg/mindmeld/internal/lobbyview"
	"github.com/jgoodhcg/mindmeld/internal/ws"
	triviatmpl "github.com/jgoodhcg/mindmeld/templates/trivia"
)

// TriviaGame implements the Game interface for the trivia game type.
type TriviaGame struct {
	queries  *db.Queries
	dbPool   *pgxpool.Pool
	eventBus events.Bus
	hub      *ws.Hub
}

// New creates a new TriviaGame.
func New(queries *db.Queries, dbPool *pgxpool.Pool, eventBus events.Bus, hub *ws.Hub) *TriviaGame {
	return &TriviaGame{
		queries:  queries,
		dbPool:   dbPool,
		eventBus: eventBus,
		hub:      hub,
	}
}

// Info returns metadata about the trivia game.
func (g *TriviaGame) Info() games.GameInfo {
	return games.GameInfo{
		Slug:        "trivia",
		Name:        "TRIVIA",
		Description: "Questions from minds you know",
		Ready:       true,
	}
}

// RenderContent builds the trivia game content component from lobby state.
func (g *TriviaGame) RenderContent(ctx context.Context, lobby db.Lobby, players []db.GetLobbyPlayersRow, player db.Player, isHost bool) templ.Component {
	var activeRound db.TriviaRound
	var hasSubmitted bool
	var currentQuestion db.TriviaQuestion
	var questionActive bool
	var isAuthor bool
	var hasAnswered bool
	var submittedCount int
	var scoreboard []db.GetLobbyScoreboardRow
	var roundScoreboard []db.GetRoundScoreboardRow
	var distribution []events.AnswerStat
	var totalAnswers int
	var totalExpectedAnswers int
	var submissionExpectedCount int
	var hostTransferOptions []lobbyview.HostTransferOption
	now := time.Now()

	if isHost {
		hostTransferOptions = lobbyview.BuildHostTransferOptions(players, player.ID.String(), func(playerID string) bool {
			return g.hub.Presence(lobby.Code, playerID).IsConnected()
		})
	}

	if lobby.Phase == "playing" {
		var err error
		activeRound, err = g.queries.GetActiveRound(ctx, lobby.ID)
		if err != nil {
			log.Printf("Error fetching active round: %v", err)
		} else {
			if activeRound.Phase == "submitting" {
				questions, err := g.queries.GetQuestionsForRound(ctx, activeRound.ID)
				if err == nil {
					submittedCount = len(questions)
					submissionExpectedCount = g.countActivePlayers(lobby.Code, players, now, "")
					for _, q := range questions {
						if q.Author == player.ID {
							hasSubmitted = true
						}
					}
				}
			} else if activeRound.Phase == "playing" {
				if activeRound.CurrentQuestionID.Valid {
					var qID pgtype.UUID = activeRound.CurrentQuestionID
					questions, err := g.queries.GetQuestionsForRound(ctx, activeRound.ID)
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
						if currentQuestion.Author == player.ID {
							isAuthor = true
						}

						for _, lobbyPlayer := range players {
							if lobbyPlayer.PlayerID != currentQuestion.Author && g.hub.Presence(lobby.Code, lobbyPlayer.PlayerID.String()).IsActiveAt(now) {
								totalExpectedAnswers++
							}
						}

						answers, err := g.queries.GetAnswersForQuestion(ctx, currentQuestion.ID)
						if err == nil {
							totalAnswers = len(answers)
							for _, a := range answers {
								if a.PlayerID == player.ID {
									hasAnswered = true
								}
							}
							distribution = buildAnswerDistributionFromAnswers(currentQuestion, answers)
						}
					}
				}
			} else if activeRound.Phase == "finished" {
				scoreboard, _ = g.queries.GetLobbyScoreboard(ctx, lobby.ID)
				roundScoreboard, _ = g.queries.GetRoundScoreboard(ctx, activeRound.ID)
			}
		}
	}

	return triviatmpl.GameContent(lobby, players, activeRound, hasSubmitted, currentQuestion, questionActive, isAuthor, hasAnswered, submittedCount, submissionExpectedCount, isHost, hostTransferOptions, scoreboard, roundScoreboard, distribution, totalAnswers, totalExpectedAnswers)
}

func (g *TriviaGame) countActivePlayers(lobbyCode string, players []db.GetLobbyPlayersRow, now time.Time, excludedPlayerID string) int {
	count := 0
	for _, lobbyPlayer := range players {
		if excludedPlayerID != "" && lobbyPlayer.PlayerID.String() == excludedPlayerID {
			continue
		}
		if g.hub.Presence(lobbyCode, lobbyPlayer.PlayerID.String()).IsActiveAt(now) {
			count++
		}
	}
	return count
}

// RegisterRoutes registers trivia-specific HTTP routes.
func (g *TriviaGame) RegisterRoutes(r chi.Router) {
	r.Post("/start", g.handleStartGame)
	r.Get("/question-templates", g.handleGetQuestionTemplates)
	r.Post("/generate-question", g.handleGenerateQuestion)
	r.Post("/questions", g.handleSubmitQuestion)
	r.Post("/advance", g.handleAdvanceRound)
	r.Post("/next-question", g.handleNextQuestion)
	r.Post("/play-again", g.handlePlayAgain)
	r.Post("/answers", g.handleSubmitAnswer)
}
