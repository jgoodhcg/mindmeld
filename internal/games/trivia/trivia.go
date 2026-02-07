package trivia

import (
	"context"
	"log"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/games"
	triviatmpl "github.com/jgoodhcg/mindmeld/templates/trivia"
)

// TriviaGame implements the Game interface for the trivia game type.
type TriviaGame struct {
	queries  *db.Queries
	dbPool   *pgxpool.Pool
	eventBus events.Bus
}

// New creates a new TriviaGame.
func New(queries *db.Queries, dbPool *pgxpool.Pool, eventBus events.Bus) *TriviaGame {
	return &TriviaGame{
		queries:  queries,
		dbPool:   dbPool,
		eventBus: eventBus,
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

						answers, err := g.queries.GetAnswersForQuestion(ctx, currentQuestion.ID)
						if err == nil {
							totalAnswers = len(answers)
							for _, a := range answers {
								if a.PlayerID == player.ID {
									hasAnswered = true
								}
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
				}
			} else if activeRound.Phase == "finished" {
				scoreboard, _ = g.queries.GetLobbyScoreboard(ctx, lobby.ID)
				roundScoreboard, _ = g.queries.GetRoundScoreboard(ctx, activeRound.ID)
			}
		}
	}

	return triviatmpl.GameContent(lobby, players, activeRound, hasSubmitted, currentQuestion, questionActive, isAuthor, hasAnswered, submittedCount, isHost, scoreboard, roundScoreboard, distribution, totalAnswers)
}

// RegisterRoutes registers trivia-specific HTTP routes.
func (g *TriviaGame) RegisterRoutes(r chi.Router) {
	r.Post("/start", g.handleStartGame)
	r.Get("/question-templates", g.handleGetQuestionTemplates)
	r.Post("/questions", g.handleSubmitQuestion)
	r.Post("/advance", g.handleAdvanceRound)
	r.Post("/next-question", g.handleNextQuestion)
	r.Post("/play-again", g.handlePlayAgain)
	r.Post("/answers", g.handleSubmitAnswer)
}
