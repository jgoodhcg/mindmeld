package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

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
			}
		}

		// 2. If no valid player found, create one
		if !player.ID.Valid {
			token := uuid.NewString()

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
		ctx = auth.WithPlayer(ctx, player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
