package auth

import (
	"context"

	"github.com/jgoodhcg/mindmeld/internal/db"
)

type contextKey string

const PlayerContextKey contextKey = "player"

// WithPlayer stores a player in the context.
func WithPlayer(ctx context.Context, player db.Player) context.Context {
	return context.WithValue(ctx, PlayerContextKey, player)
}

// GetPlayer retrieves the player object from the request context.
// Since middleware ensures it's always there, this effectively always returns a player.
func GetPlayer(ctx context.Context) db.Player {
	player, ok := ctx.Value(PlayerContextKey).(db.Player)
	if !ok {
		return db.Player{}
	}
	return player
}
