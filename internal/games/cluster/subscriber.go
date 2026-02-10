package cluster

import (
	"context"
	"log"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/ws"
)

// HandleEvent processes Cluster-specific events.
func (g *ClusterGame) HandleEvent(ctx context.Context, event events.Event, hub *ws.Hub, _ *db.Queries) bool {
	switch event.Type {
	case events.EventClusterRoundStarted,
		events.EventClusterSubmissionUpdated,
		events.EventClusterRoundRevealed,
		events.EventClusterExhausted:
		ws.BroadcastUpdateTrigger(ctx, event.LobbyCode, hub)
		return true
	default:
		log.Printf("[cluster-subscriber] unhandled event type: %s", event.Type)
		return false
	}
}
