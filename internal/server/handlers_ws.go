package server

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
)

// handleWebSocket upgrades the HTTP connection to a WebSocket
// and registers it with the hub for the given lobby.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	// Verify the lobby exists
	_, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	// Accept the WebSocket connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Allow connections from any origin for development
		// In production, you may want to restrict this
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("[ws] Failed to accept connection: %v", err)
		return
	}

	// Register with hub
	s.hub.Register(code, conn)

	// Ensure cleanup on exit
	defer func() {
		s.hub.Unregister(code, conn)
		conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	// Keep the connection alive by reading messages
	// We don't expect any messages from the client for now,
	// but we need to read to detect disconnection
	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			// Connection closed or error
			log.Printf("[ws] Connection closed for lobby %s: %v", code, err)
			return
		}
	}
}
