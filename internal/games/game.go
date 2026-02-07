package games

import (
	"context"
	"net/http"
	"sync"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/ws"
)

// GameInfo describes a game for display on the platform page.
type GameInfo struct {
	Slug        string
	Name        string
	Description string
	Ready       bool
}

// Game defines the interface that all game types must implement.
type Game interface {
	Info() GameInfo
	RegisterRoutes(r chi.Router)
	RenderContent(ctx context.Context, lobby db.Lobby, players []db.GetLobbyPlayersRow, player db.Player, isHost bool) templ.Component
	HandleEvent(ctx context.Context, event events.Event, hub *ws.Hub, queries *db.Queries) bool
}

// Registry holds all registered games.
type Registry struct {
	mu    sync.RWMutex
	games map[string]Game
	// placeholders stores game info for games not yet implemented
	placeholders []GameInfo
}

// NewRegistry creates a new game registry.
func NewRegistry() *Registry {
	return &Registry{
		games: make(map[string]Game),
	}
}

// Register adds a game to the registry.
func (r *Registry) Register(game Game) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.games[game.Info().Slug] = game
}

// Get returns a game by slug.
func (r *Registry) Get(slug string) (Game, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.games[slug]
	return g, ok
}

// Each iterates over all registered games.
func (r *Registry) Each(fn func(slug string, game Game)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for slug, game := range r.games {
		fn(slug, game)
	}
}

// RegisterPlaceholder adds a coming-soon game to the registry.
func (r *Registry) RegisterPlaceholder(info GameInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.placeholders = append(r.placeholders, info)
}

// GetHandler returns a game as a ws.GameHandler interface.
// This satisfies the ws.GameRegistry interface.
func (r *Registry) GetHandler(gameType string) (ws.GameHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.games[gameType]
	if !ok {
		return nil, false
	}
	return g, true
}

// AllInfo returns GameInfo for all games (registered + placeholders).
func (r *Registry) AllInfo() []GameInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var infos []GameInfo
	for _, game := range r.games {
		infos = append(infos, game.Info())
	}
	infos = append(infos, r.placeholders...)
	return infos
}

// HandlerFunc is a helper type for game handlers that need access to the request.
type HandlerFunc func(http.ResponseWriter, *http.Request)
