// Package ws provides WebSocket connection management for real-time updates.
package ws

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type playerState struct {
	connections    int
	disconnectedAt time.Time
	graceTimer     *time.Timer
}

// Hub manages WebSocket connections grouped by lobby.
type Hub struct {
	// lobbies maps lobby codes to their connected clients with player IDs (as UUID strings)
	lobbies map[string]map[*websocket.Conn]string
	players map[string]map[string]*playerState
	mu      sync.RWMutex

	disconnectGracePeriod time.Duration
	presenceHandler       func(lobbyCode string, update PresenceUpdate)
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		lobbies:               make(map[string]map[*websocket.Conn]string),
		players:               make(map[string]map[string]*playerState),
		disconnectGracePeriod: defaultDisconnectGracePeriod,
	}
}

// SetPresenceHandler registers a callback for presence changes.
func (h *Hub) SetPresenceHandler(handler func(lobbyCode string, update PresenceUpdate)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.presenceHandler = handler
}

// SetDisconnectGracePeriod overrides the default disconnect grace period.
// Intended for tests.
func (h *Hub) SetDisconnectGracePeriod(gracePeriod time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.disconnectGracePeriod = gracePeriod
}

func (h *Hub) getOrCreatePlayerState(lobbyCode string, playerID string) *playerState {
	if h.players[lobbyCode] == nil {
		h.players[lobbyCode] = make(map[string]*playerState)
	}
	state, ok := h.players[lobbyCode][playerID]
	if !ok {
		state = &playerState{}
		h.players[lobbyCode][playerID] = state
	}
	return state
}

func (h *Hub) notifyPresence(lobbyCode string, update PresenceUpdate) {
	h.mu.RLock()
	handler := h.presenceHandler
	h.mu.RUnlock()
	if handler != nil {
		handler(lobbyCode, update)
	}
}

// Register adds a connection to a lobby with the associated player ID (UUID string).
func (h *Hub) Register(lobbyCode string, conn *websocket.Conn, playerID string) {
	shouldNotify := false
	totalConnections := 0

	h.mu.Lock()

	if h.lobbies[lobbyCode] == nil {
		h.lobbies[lobbyCode] = make(map[*websocket.Conn]string)
	}
	h.lobbies[lobbyCode][conn] = playerID
	state := h.getOrCreatePlayerState(lobbyCode, playerID)
	if state.connections == 0 {
		shouldNotify = true
	}
	state.connections++
	state.disconnectedAt = time.Time{}
	if state.graceTimer != nil {
		state.graceTimer.Stop()
		state.graceTimer = nil
	}
	totalConnections = len(h.lobbies[lobbyCode])
	h.mu.Unlock()

	log.Printf("[ws] Client connected to lobby %s (playerID: %s, %d total)", lobbyCode, playerID, totalConnections)
	if shouldNotify {
		h.notifyPresence(lobbyCode, PresenceUpdate{PlayerID: playerID, Connected: true})
	}
}

// Unregister removes a connection from a lobby.
func (h *Hub) Unregister(lobbyCode string, conn *websocket.Conn) {
	var (
		playerID     string
		shouldNotify bool
	)

	h.mu.Lock()

	if clients, ok := h.lobbies[lobbyCode]; ok {
		playerID = clients[conn]
		delete(clients, conn)
		log.Printf("[ws] Client disconnected from lobby %s (%d remaining)", lobbyCode, len(clients))

		if playerID != "" {
			state := h.getOrCreatePlayerState(lobbyCode, playerID)
			if state.connections > 0 {
				state.connections--
			}
			if state.connections == 0 {
				shouldNotify = true
				state.disconnectedAt = time.Now()
				if state.graceTimer != nil {
					state.graceTimer.Stop()
				}
				gracePeriod := h.disconnectGracePeriod
				playerIDCopy := playerID
				state.graceTimer = time.AfterFunc(gracePeriod, func() {
					h.handleGraceExpiry(lobbyCode, playerIDCopy)
				})
			}
		}

		// Clean up empty lobbies
		if len(clients) == 0 {
			delete(h.lobbies, lobbyCode)
		}
	}
	h.mu.Unlock()

	if shouldNotify {
		h.notifyPresence(lobbyCode, PresenceUpdate{PlayerID: playerID, Connected: false})
	}
}

func (h *Hub) handleGraceExpiry(lobbyCode string, playerID string) {
	h.mu.Lock()
	players := h.players[lobbyCode]
	state, ok := players[playerID]
	if !ok {
		h.mu.Unlock()
		return
	}
	if state.connections > 0 || state.disconnectedAt.IsZero() {
		h.mu.Unlock()
		return
	}
	state.graceTimer = nil
	h.mu.Unlock()

	h.notifyPresence(lobbyCode, PresenceUpdate{PlayerID: playerID, Connected: false, GraceExpired: true})
}

// Broadcast sends a message to all connections in a lobby.
func (h *Hub) Broadcast(ctx context.Context, lobbyCode string, message []byte) {
	h.mu.RLock()
	clients := make([]*websocket.Conn, 0)
	if lobby, ok := h.lobbies[lobbyCode]; ok {
		for conn := range lobby {
			clients = append(clients, conn)
		}
	}
	h.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	log.Printf("[ws] Broadcasting to %d clients in lobby %s", len(clients), lobbyCode)

	for _, conn := range clients {
		err := conn.Write(ctx, websocket.MessageText, message)
		if err != nil {
			log.Printf("[ws] Error writing to client: %v", err)
			// Connection will be cleaned up by the read loop
		}
	}
}

// ConnectionCount returns the number of connections in a lobby.
func (h *Hub) ConnectionCount(lobbyCode string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.lobbies[lobbyCode]; ok {
		return len(clients)
	}
	return 0
}

// LiveLobbyCount returns the number of lobbies with active connections.
func (h *Hub) LiveLobbyCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.lobbies)
}

// Presence returns the current presence snapshot for a player in a lobby.
func (h *Hub) Presence(lobbyCode string, playerID string) PlayerPresence {
	h.mu.RLock()
	defer h.mu.RUnlock()

	state, ok := h.players[lobbyCode][playerID]
	if !ok {
		return PlayerPresence{GracePeriod: h.disconnectGracePeriod}
	}

	return PlayerPresence{
		ConnectionCount: state.connections,
		DisconnectedAt:  state.disconnectedAt,
		GracePeriod:     h.disconnectGracePeriod,
	}
}

// Snapshot returns all known player presence for a lobby.
func (h *Hub) Snapshot(lobbyCode string) map[string]PlayerPresence {
	h.mu.RLock()
	defer h.mu.RUnlock()

	snapshot := make(map[string]PlayerPresence)
	for playerID, state := range h.players[lobbyCode] {
		snapshot[playerID] = PlayerPresence{
			ConnectionCount: state.connections,
			DisconnectedAt:  state.disconnectedAt,
			GracePeriod:     h.disconnectGracePeriod,
		}
	}
	return snapshot
}

// BroadcastPersonalized sends a personalized message to each connection in a lobby.
// The renderFunc is called for each player ID (UUID string) and should return the message bytes for that player.
func (h *Hub) BroadcastPersonalized(ctx context.Context, lobbyCode string, renderFunc func(playerID string) []byte) {
	h.mu.RLock()
	type clientInfo struct {
		conn     *websocket.Conn
		playerID string
	}
	clients := make([]clientInfo, 0)
	if lobby, ok := h.lobbies[lobbyCode]; ok {
		for conn, playerID := range lobby {
			clients = append(clients, clientInfo{conn: conn, playerID: playerID})
		}
	}
	h.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	log.Printf("[ws] Broadcasting personalized messages to %d clients in lobby %s", len(clients), lobbyCode)

	for _, client := range clients {
		message := renderFunc(client.playerID)
		if message != nil {
			err := client.conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				log.Printf("[ws] Error writing to client (playerID: %s): %v", client.playerID, err)
			}
		}
	}
}
