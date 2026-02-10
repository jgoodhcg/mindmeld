package cluster

import (
	"bytes"
	"context"
	"log"

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
