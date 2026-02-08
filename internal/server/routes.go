package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jgoodhcg/mindmeld/internal/games"
)

func (s *Server) routes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(s.playerIdentityMiddleware)

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// App Routes
	s.router.Get("/", s.handlePlatform)
	s.router.Get("/trivia", s.handleTriviaHome)
	s.router.Post("/trivia/join", s.handleJoinByCode)
	s.router.Post("/lobbies", s.handleCreateLobby)
	s.router.Get("/lobbies/{code}", s.handleLobbyRoom)
	s.router.Get("/lobbies/{code}/content", s.handleGetGameContent)
	s.router.Post("/lobbies/{code}/join", s.handleJoinLobby)

	// WebSocket for real-time updates
	s.router.Get("/lobbies/{code}/ws", s.handleWebSocket)

	// Dynamic game routes: /lobbies/{code}/{game_slug}/...
	// These paths are action namespaces, not the canonical lobby URL.
	// Every game-scoped action is guarded so {game_slug} must match the
	// lobby's current game_type before the game handler executes.
	s.games.Each(func(_ string, game games.Game) {
		slug := game.Info().Slug
		s.router.Route("/lobbies/{code}/"+slug, func(r chi.Router) {
			r.Use(s.requireLobbyGameType(slug))
			game.RegisterRoutes(r)
		})
	})

	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}
