package trivia

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/ws"
	triviatmpl "github.com/jgoodhcg/mindmeld/templates/trivia"
)

// HandleEvent processes trivia-specific WebSocket events.
// Returns true if the event was handled.
func (g *TriviaGame) HandleEvent(ctx context.Context, event events.Event, hub *ws.Hub, queries *db.Queries) bool {
	switch event.Type {
	case events.EventGameStarted:
		g.broadcastGameStarted(ctx, event.LobbyCode, event.Payload.(events.GameStartedPayload), hub, queries)
		return true
	case events.EventQuestionSubmitted:
		g.broadcastQuestionSubmitted(ctx, event.LobbyCode, event.Payload.(events.QuestionSubmittedPayload), hub)
		return true
	case events.EventRoundAdvanced:
		g.broadcastRoundAdvanced(ctx, event.LobbyCode, event.Payload.(events.RoundAdvancedPayload), hub)
		return true
	case events.EventQuestionRevealed:
		g.broadcastQuestionRevealed(ctx, event.LobbyCode, event.Payload.(events.QuestionRevealedPayload), hub)
		return true
	case events.EventAnswerSubmitted:
		g.broadcastAnswerSubmitted(ctx, event.LobbyCode, event.Payload.(events.AnswerSubmittedPayload), hub, queries)
		return true
	case events.EventNewRoundCreated:
		g.broadcastNewRoundCreated(ctx, event.LobbyCode, event.Payload.(events.NewRoundCreatedPayload), hub)
		return true
	case events.EventPlayerPresence:
		payload, ok := event.Payload.(events.PlayerPresencePayload)
		if !ok || !payload.GraceExpired {
			return false
		}
		return g.handleGraceExpired(ctx, event.LobbyCode)
	default:
		return false
	}
}

func (g *TriviaGame) broadcastGameStarted(ctx context.Context, lobbyCode string, payload events.GameStartedPayload, hub *ws.Hub, queries *db.Queries) {
	lobby, err := queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[trivia-subscriber] Failed to get lobby %s: %v", lobbyCode, err)
		return
	}

	round, err := queries.GetActiveRound(ctx, lobby.ID)
	if err != nil {
		log.Printf("[trivia-subscriber] Failed to get active round for lobby %s: %v", lobbyCode, err)
		return
	}

	players, err := queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[trivia-subscriber] Failed to get players for lobby %s: %v", lobbyCode, err)
		return
	}

	var buf bytes.Buffer
	var emptyQuestion db.TriviaQuestion
	var emptyScoreboard []db.GetLobbyScoreboardRow
	var emptyRoundScoreboard []db.GetRoundScoreboardRow
	var emptyDistribution []events.AnswerStat
	err = triviatmpl.GameContent(
		lobby,
		players,
		round,
		false,
		emptyQuestion,
		false,
		false,
		false,
		0,
		g.countActivePlayers(lobbyCode, players, time.Now(), ""),
		false,
		emptyScoreboard,
		emptyRoundScoreboard,
		emptyDistribution,
		0,
		0,
	).Render(ctx, &buf)
	if err != nil {
		log.Printf("[trivia-subscriber] Failed to render game content: %v", err)
		return
	}

	hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}

func (g *TriviaGame) broadcastQuestionSubmitted(ctx context.Context, lobbyCode string, payload events.QuestionSubmittedPayload, hub *ws.Hub) {
	hub.BroadcastPersonalized(ctx, lobbyCode, func(playerID string) []byte {
		isHost := playerID == payload.HostPlayerID

		var buf bytes.Buffer
		err := triviatmpl.SubmitStatus(
			payload.SubmittedCount,
			payload.TotalPlayers,
			lobbyCode,
			isHost,
			true,
		).Render(ctx, &buf)
		if err != nil {
			log.Printf("[trivia-subscriber] Failed to render submit status for player %s: %v", playerID, err)
			return nil
		}
		return buf.Bytes()
	})
}

func (g *TriviaGame) broadcastRoundAdvanced(ctx context.Context, lobbyCode string, payload events.RoundAdvancedPayload, hub *ws.Hub) {
	log.Printf("[trivia-subscriber] Round advanced for lobby %s, round %d", lobbyCode, payload.RoundNumber)
	ws.BroadcastUpdateTrigger(ctx, lobbyCode, hub)
}

func (g *TriviaGame) broadcastQuestionRevealed(ctx context.Context, lobbyCode string, payload events.QuestionRevealedPayload, hub *ws.Hub) {
	log.Printf("[trivia-subscriber] Question revealed for lobby %s", lobbyCode)
	ws.BroadcastUpdateTrigger(ctx, lobbyCode, hub)
}

func (g *TriviaGame) broadcastAnswerSubmitted(ctx context.Context, lobbyCode string, payload events.AnswerSubmittedPayload, hub *ws.Hub, queries *db.Queries) {
	if payload.QuestionComplete {
		return
	}

	lobby, err := queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[trivia-subscriber] Error fetching lobby: %v", err)
		return
	}

	activeRound, err := queries.GetActiveRound(ctx, lobby.ID)
	if err != nil || !activeRound.CurrentQuestionID.Valid {
		log.Printf("[trivia-subscriber] Error fetching active round/question: %v", err)
		return
	}

	roundQuestions, err := queries.GetQuestionsForRound(ctx, activeRound.ID)
	var currentQuestion db.TriviaQuestion
	found := false
	if err == nil {
		for _, q := range roundQuestions {
			if q.ID == activeRound.CurrentQuestionID {
				currentQuestion = q
				found = true
				break
			}
		}
	}
	if !found {
		return
	}

	answers, err := queries.GetAnswersForQuestion(ctx, currentQuestion.ID)
	if err != nil {
		return
	}
	answeredMap := make(map[string]bool)
	totalAnswers := len(answers)
	for _, a := range answers {
		answeredMap[a.PlayerID.String()] = true
	}

	hub.BroadcastPersonalized(ctx, lobbyCode, func(playerID string) []byte {
		hasAnswered := answeredMap[playerID]
		isAuthor := currentQuestion.Author.String() == playerID

		if hasAnswered || isAuthor {
			var buf bytes.Buffer
			err := triviatmpl.QuestionResults(
				lobbyCode,
				currentQuestion,
				payload.Distribution,
				totalAnswers,
				payload.TotalExpected,
				false,
				false,
			).Render(ctx, &buf)

			if err != nil {
				log.Printf("[trivia-subscriber] Error rendering live graph: %v", err)
				return nil
			}
			return buf.Bytes()
		}

		var buf2 bytes.Buffer
		err := triviatmpl.AnswerStatus(payload.AnsweredCount, payload.TotalExpected, true).Render(ctx, &buf2)
		if err != nil {
			log.Printf("[trivia-subscriber] Error rendering answer status: %v", err)
			return nil
		}
		return buf2.Bytes()
	})
}

func (g *TriviaGame) broadcastNewRoundCreated(ctx context.Context, lobbyCode string, payload events.NewRoundCreatedPayload, hub *ws.Hub) {
	log.Printf("[trivia-subscriber] New round created for lobby %s, round %d", lobbyCode, payload.RoundNumber)
	ws.BroadcastUpdateTrigger(ctx, lobbyCode, hub)
}

func (g *TriviaGame) handleGraceExpired(ctx context.Context, lobbyCode string) bool {
	lobby, err := g.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil || lobby.Phase != "playing" {
		return false
	}

	round, err := g.queries.GetActiveRound(ctx, lobby.ID)
	if err != nil || round.Phase != "playing" || round.QuestionState != "answering" || !round.CurrentQuestionID.Valid {
		return false
	}

	questions, err := g.queries.GetQuestionsForRound(ctx, round.ID)
	if err != nil {
		return false
	}

	var currentQuestion db.TriviaQuestion
	found := false
	for _, question := range questions {
		if question.ID == round.CurrentQuestionID {
			currentQuestion = question
			found = true
			break
		}
	}
	if !found {
		return false
	}

	players, err := g.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		return false
	}
	expected := g.countActivePlayers(lobbyCode, players, time.Now(), currentQuestion.Author.String())
	currentCount, err := g.queries.CountAnswersForQuestion(ctx, currentQuestion.ID)
	if err != nil || int(currentCount) < expected {
		return false
	}

	rawStats, err := g.queries.GetAnswerStats(ctx, currentQuestion.ID)
	if err != nil {
		return false
	}
	distribution := make([]events.AnswerStat, 0, len(rawStats))
	for _, stat := range rawStats {
		distribution = append(distribution, events.AnswerStat{Answer: stat.SelectedAnswer, Count: int(stat.Count)})
	}

	if err = g.queries.UpdateRoundQuestionState(ctx, db.UpdateRoundQuestionStateParams{
		ID:                round.ID,
		CurrentQuestionID: currentQuestion.ID,
		QuestionState:     "revealed",
	}); err != nil {
		return false
	}

	g.eventBus.Publish(ctx, events.Event{
		Type:      events.EventQuestionRevealed,
		LobbyCode: lobbyCode,
		Payload: events.QuestionRevealedPayload{
			QuestionID:   currentQuestion.ID.String(),
			Distribution: distribution,
		},
	})
	return true
}
