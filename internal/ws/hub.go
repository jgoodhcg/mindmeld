// Package ws provides WebSocket connection management for real-time updates.
package ws

import (
	"context"
	"log"
	"sync"

	"github.com/coder/websocket"
)

// Hub manages WebSocket connections grouped by lobby.
type Hub struct {
	// lobbies maps lobby codes to their connected clients with player IDs (as UUID strings)
	lobbies map[string]map[*websocket.Conn]string
	mu      sync.RWMutex
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		lobbies: make(map[string]map[*websocket.Conn]string),
	}
}

// Register adds a connection to a lobby with the associated player ID (UUID string).
func (h *Hub) Register(lobbyCode string, conn *websocket.Conn, playerID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lobbies[lobbyCode] == nil {
		h.lobbies[lobbyCode] = make(map[*websocket.Conn]string)
	}
	h.lobbies[lobbyCode][conn] = playerID

	log.Printf("[ws] Client connected to lobby %s (playerID: %s, %d total)", lobbyCode, playerID, len(h.lobbies[lobbyCode]))
}

// Unregister removes a connection from a lobby.
func (h *Hub) Unregister(lobbyCode string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.lobbies[lobbyCode]; ok {
		delete(clients, conn)
		log.Printf("[ws] Client disconnected from lobby %s (%d remaining)", lobbyCode, len(clients))

		// Clean up empty lobbies
		if len(clients) == 0 {
			delete(h.lobbies, lobbyCode)
		}
	}
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
