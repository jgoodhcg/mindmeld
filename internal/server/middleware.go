package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

type contextKey string

const playerContextKey contextKey = "player"

// playerIdentityMiddleware ensures every request has a valid player identity associated with it.
// It checks for a device_token cookie. If missing or invalid, it creates a new player.
func (s *Server) playerIdentityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var player db.Player
		var err error

		// 1. Try to get token from cookie
		cookie, err := r.Cookie("device_token")
		if err == nil {
			// Token exists, verify in DB
			player, err = s.queries.GetPlayerByDeviceToken(ctx, cookie.Value)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Error fetching player by token: %v", err)
				// Continue without player? Or fail?
				// For now, if DB error, we treat as new session to be safe or fail 500.
				// Let's treat as new session (fall through) unless it's a real DB outage.
			}
		}

		// 2. If no valid player found, create one
		if !player.ID.Valid {
			token := uuid.NewString()

			// Create new player in DB
			// We pass sql.NullUUID{} for user_id as it's optional
			player, err = s.queries.CreatePlayer(ctx, db.CreatePlayerParams{
				DeviceToken: token,
			})
			if err != nil {
				log.Printf("Failed to create new player: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set long-lived cookie (1 year)
			http.SetCookie(w, &http.Cookie{
				Name:     "device_token",
				Value:    token,
				Path:     "/",
				Expires:  time.Now().Add(365 * 24 * time.Hour),
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
		}

		// 3. Store player in context
		ctx = context.WithValue(ctx, playerContextKey, player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPlayer retrieves the player object from the request context.
// Since middleware ensures it's always there, this effectively always returns a player.
func GetPlayer(ctx context.Context) db.Player {
	player, ok := ctx.Value(playerContextKey).(db.Player)
	if !ok {
		// Should technically never happen if middleware is applied
		return db.Player{}
	}
	return player
}
