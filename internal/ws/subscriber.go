package ws

import (
	"bytes"
	"context"
	"log"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/templates"
)

// Subscriber listens to events and broadcasts updates to WebSocket clients.
type Subscriber struct {
	hub     *Hub
	queries *db.Queries
}

// NewSubscriber creates a new event subscriber.
func NewSubscriber(hub *Hub, queries *db.Queries) *Subscriber {
	return &Subscriber{
		hub:     hub,
		queries: queries,
	}
}

// HandleEvent processes an event and broadcasts appropriate updates.
func (s *Subscriber) HandleEvent(ctx context.Context, event events.Event) {
	switch event.Type {
	case events.EventPlayerJoined, events.EventPlayerLeft:
		s.broadcastPlayerList(ctx, event.LobbyCode)
	case events.EventGameStarted:
		s.broadcastGameStarted(ctx, event.LobbyCode, event.Payload.(events.GameStartedPayload))
	case events.EventQuestionSubmitted:
		s.broadcastQuestionSubmitted(ctx, event.LobbyCode, event.Payload.(events.QuestionSubmittedPayload))
	case events.EventRoundAdvanced:
		s.broadcastRoundAdvanced(ctx, event.LobbyCode, event.Payload.(events.RoundAdvancedPayload))
	case events.EventQuestionRevealed:
		s.broadcastQuestionRevealed(ctx, event.LobbyCode, event.Payload.(events.QuestionRevealedPayload))
	case events.EventAnswerSubmitted:
		s.broadcastAnswerSubmitted(ctx, event.LobbyCode, event.Payload.(events.AnswerSubmittedPayload))
	case events.EventNewRoundCreated:
		s.broadcastNewRoundCreated(ctx, event.LobbyCode, event.Payload.(events.NewRoundCreatedPayload))
	default:
		log.Printf("[ws-subscriber] Unhandled event type: %s", event.Type)
	}
}

// broadcastPlayerList fetches the current player list and broadcasts it.
func (s *Subscriber) broadcastPlayerList(ctx context.Context, lobbyCode string) {
	// Get the lobby to find its ID
	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get lobby %s: %v", lobbyCode, err)
		return
	}

	// Fetch current players
	players, err := s.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get players for lobby %s: %v", lobbyCode, err)
		return
	}

	// Render the player list partial to HTML
	var buf bytes.Buffer
	err = templates.PlayerList(players).Render(ctx, &buf)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to render player list: %v", err)
		return
	}

	// Broadcast to all connected clients in this lobby
	s.hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}

// broadcastGameStarted broadcasts the game content when the game starts.
// All players see the SubmitQuestion form since no one has submitted yet.
func (s *Subscriber) broadcastGameStarted(ctx context.Context, lobbyCode string, payload events.GameStartedPayload) {
	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get lobby %s: %v", lobbyCode, err)
		return
	}

	// Get the active round
	round, err := s.queries.GetActiveRound(ctx, lobby.ID)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get active round for lobby %s: %v", lobbyCode, err)
		return
	}

	// Get players for the player count
	players, err := s.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to get players for lobby %s: %v", lobbyCode, err)
		return
	}

	// Render the GameContent partial with initial state (no one has submitted)
	// Using zero values for question-related params since we're in submitting phase
	var buf bytes.Buffer
	var emptyQuestion db.TriviaQuestion
	var emptyScoreboard []db.GetLobbyScoreboardRow
	var emptyRoundScoreboard []db.GetRoundScoreboardRow
	var emptyDistribution []events.AnswerStat
	err = templates.GameContent(
		lobby,
		players,
		round,
		false, // hasSubmitted - no one has submitted yet
		emptyQuestion,
		false, // questionActive
		false, // isAuthor
		false, // hasAnswered
		0,     // submittedCount
		false, // isHost - this is tricky, but the form doesn't need it in submitting state
		emptyScoreboard,
		emptyRoundScoreboard,
		emptyDistribution,
		0,
	).Render(ctx, &buf)
	if err != nil {
		log.Printf("[ws-subscriber] Failed to render game content: %v", err)
		return
	}

	s.hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}

// broadcastQuestionSubmitted broadcasts the updated submit status counter.
// Uses personalized broadcasting so the host sees the start button.
func (s *Subscriber) broadcastQuestionSubmitted(ctx context.Context, lobbyCode string, payload events.QuestionSubmittedPayload) {
	s.hub.BroadcastPersonalized(ctx, lobbyCode, func(playerID string) []byte {
		isHost := playerID == payload.HostPlayerID

		var buf bytes.Buffer
		err := templates.SubmitStatus(
			payload.SubmittedCount,
			payload.TotalPlayers,
			lobbyCode,
			isHost,
		).Render(ctx, &buf)
		if err != nil {
			log.Printf("[ws-subscriber] Failed to render submit status for player %s: %v", playerID, err)
			return nil
		}
		return buf.Bytes()
	})
}

// broadcastRoundAdvanced broadcasts when the round advances to playing phase.
// Since each player sees personalized content (author vs answerer), we trigger a page refresh.
func (s *Subscriber) broadcastRoundAdvanced(ctx context.Context, lobbyCode string, payload events.RoundAdvancedPayload) {
	log.Printf("[ws-subscriber] Round advanced for lobby %s, round %d", lobbyCode, payload.RoundNumber)

	// Broadcast a script that triggers a page reload for all clients
	// This ensures each player fetches their personalized view
	refreshHTML := []byte(`<div id="game-content"><script>window.location.reload()</script></div>`)
	s.hub.Broadcast(ctx, lobbyCode, refreshHTML)
}

// broadcastQuestionRevealed broadcasts when everyone has answered.
// Triggers a page refresh so everyone sees the revealed results.
func (s *Subscriber) broadcastQuestionRevealed(ctx context.Context, lobbyCode string, payload events.QuestionRevealedPayload) {
	log.Printf("[ws-subscriber] Question revealed for lobby %s", lobbyCode)
	refreshHTML := []byte(`<div id="game-content"><script>window.location.reload()</script></div>`)
	s.hub.Broadcast(ctx, lobbyCode, refreshHTML)
}

// broadcastAnswerSubmitted broadcasts answer progress updates.
// If question is complete, handled by broadcastQuestionRevealed.
// Otherwise, sends Live Graph to players who have already answered.
func (s *Subscriber) broadcastAnswerSubmitted(ctx context.Context, lobbyCode string, payload events.AnswerSubmittedPayload) {
	if payload.QuestionComplete {
		// Do nothing, QuestionRevealed event will handle the reload
		return
	}

	// Fetch necessary data to render the graph
	lobby, err := s.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil {
		log.Printf("[ws-subscriber] Error fetching lobby: %v", err)
		return
	}

	activeRound, err := s.queries.GetActiveRound(ctx, lobby.ID)
	if err != nil || !activeRound.CurrentQuestionID.Valid {
		log.Printf("[ws-subscriber] Error fetching active round/question: %v", err)
		return
	}

	// Fetch Current Question
	questions, err := s.queries.GetQuestionsForRound(ctx, activeRound.ID)
	var currentQuestion db.TriviaQuestion
	found := false
	if err == nil {
		for _, q := range questions {
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

	// Fetch who has answered
	answers, err := s.queries.GetAnswersForQuestion(ctx, currentQuestion.ID)
	if err != nil {
		return
	}
	answeredMap := make(map[string]bool)
	totalAnswers := len(answers)
	for _, a := range answers {
		answeredMap[a.PlayerID.String()] = true
	}

	// Broadcast Personalized
	s.hub.BroadcastPersonalized(ctx, lobbyCode, func(playerID string) []byte {
		// Only send update if player has answered OR is the author
		hasAnswered := answeredMap[playerID]
		isAuthor := currentQuestion.Author.String() == playerID
		
		if hasAnswered || isAuthor {
			var buf bytes.Buffer
			// QuestionResults(lobbyCode string, question db.TriviaQuestion, distribution []events.AnswerStat, totalAnswers int, isRevealed bool, isHost bool)
			err := templates.QuestionResults(
				lobbyCode,
				currentQuestion,
				payload.Distribution,
				totalAnswers,
				false, // Not revealed yet
				false, // Host button hidden anyway
			).Render(ctx, &buf)
			
			if err != nil {
				log.Printf("[ws-subscriber] Error rendering live graph: %v", err)
				return nil
			}
			return buf.Bytes()
		}
		
		// Send answer status to players still answering
		var buf2 bytes.Buffer
		err := templates.AnswerStatus(payload.AnsweredCount, payload.TotalExpected).Render(ctx, &buf2)
		if err != nil {
			log.Printf("[ws-subscriber] Error rendering answer status: %v", err)
			return nil
		}
		return buf2.Bytes()
	})
}

// broadcastNewRoundCreated broadcasts when a new round is created (Play Again).
// Triggers a page refresh so all players see the new question submission form.
func (s *Subscriber) broadcastNewRoundCreated(ctx context.Context, lobbyCode string, payload events.NewRoundCreatedPayload) {
	log.Printf("[ws-subscriber] New round created for lobby %s, round %d", lobbyCode, payload.RoundNumber)

	// Trigger page refresh for all clients to see the new round
	refreshHTML := []byte(`<div id="game-content"><script>window.location.reload()</script></div>`)
	s.hub.Broadcast(ctx, lobbyCode, refreshHTML)
}
