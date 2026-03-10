package cluster

import (
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/ws"
	clustertmpl "github.com/jgoodhcg/mindmeld/templates/cluster"
)

// HandleEvent processes Cluster-specific events.
func (g *ClusterGame) HandleEvent(ctx context.Context, event events.Event, hub *ws.Hub, _ *db.Queries) bool {
	switch event.Type {
	case events.EventClusterRoundStarted,
		events.EventClusterRoundRevealed,
		events.EventClusterExhausted:
		ws.BroadcastUpdateTrigger(ctx, event.LobbyCode, hub)
		return true
	case events.EventClusterSubmissionUpdated:
		payload, ok := event.Payload.(events.ClusterSubmissionUpdatedPayload)
		if !ok {
			// Fallback to full refresh if payload is missing.
			ws.BroadcastUpdateTrigger(ctx, event.LobbyCode, hub)
			return true
		}
		g.broadcastSubmissionStatus(ctx, event.LobbyCode, payload, hub)
		return true
	case events.EventPlayerPresence:
		payload, ok := event.Payload.(events.PlayerPresencePayload)
		if !ok || !payload.GraceExpired {
			return false
		}
		return g.handleGraceExpired(ctx, event.LobbyCode)
	default:
		log.Printf("[cluster-subscriber] unhandled event type: %s", event.Type)
		return false
	}
}

func (g *ClusterGame) broadcastSubmissionStatus(ctx context.Context, lobbyCode string, payload events.ClusterSubmissionUpdatedPayload, hub *ws.Hub) {
	var buf bytes.Buffer
	if err := clustertmpl.SubmissionStatus(payload.SubmittedCount, payload.TotalPlayers, true).Render(ctx, &buf); err != nil {
		log.Printf("[cluster-subscriber] failed to render submission status for lobby %s: %v", lobbyCode, err)
		return
	}
	hub.Broadcast(ctx, lobbyCode, buf.Bytes())
}

func (g *ClusterGame) handleGraceExpired(ctx context.Context, lobbyCode string) bool {
	lobby, err := g.queries.GetLobbyByCode(ctx, lobbyCode)
	if err != nil || !strings.EqualFold(lobby.Phase, "playing") {
		return false
	}

	players, err := g.queries.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		return false
	}
	expectedPlayers := g.countActivePlayers(lobbyCode, players, time.Now())
	if expectedPlayers == 0 {
		return false
	}

	tx, err := g.dbPool.Begin(ctx)
	if err != nil {
		return false
	}
	defer tx.Rollback(ctx)

	round, err := g.getLatestRound(ctx, tx, lobby.ID)
	if err != nil || (round.CentroidX.Valid && round.CentroidY.Valid) {
		return false
	}

	submissionCount, err := g.countRoundSubmissions(ctx, tx, round.ID)
	if err != nil || submissionCount < expectedPlayers {
		return false
	}

	submissions, err := g.getRoundSubmissions(ctx, tx, round.ID)
	if err != nil {
		return false
	}
	points := make([]Point, 0, len(submissions))
	for _, sub := range submissions {
		points = append(points, Point{X: sub.X, Y: sub.Y})
	}
	centroidX, centroidY, ok := CalculateCentroid(points)
	if !ok {
		return false
	}
	if err = g.setRoundCentroid(ctx, tx, round.ID, centroidX, centroidY); err != nil {
		return false
	}
	if err = tx.Commit(ctx); err != nil {
		return false
	}

	g.eventBus.Publish(ctx, events.Event{Type: events.EventClusterRoundRevealed, LobbyCode: lobbyCode})
	return true
}
