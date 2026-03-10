package ws

import (
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestHubPresenceLifecycle(t *testing.T) {
	hub := NewHub()
	hub.SetDisconnectGracePeriod(20 * time.Millisecond)

	updates := make(chan PresenceUpdate, 3)
	hub.SetPresenceHandler(func(lobbyCode string, update PresenceUpdate) {
		if lobbyCode == "ABC123" {
			updates <- update
		}
	})

	conn := new(websocket.Conn)
	hub.Register("ABC123", conn, "player-1")
	assertUpdate(t, updates, PresenceUpdate{PlayerID: "player-1", Connected: true})

	hub.Unregister("ABC123", conn)
	assertUpdate(t, updates, PresenceUpdate{PlayerID: "player-1", Connected: false})
	assertUpdate(t, updates, PresenceUpdate{PlayerID: "player-1", Connected: false, GraceExpired: true})
}

func TestHubReconnectCancelsGraceExpiry(t *testing.T) {
	hub := NewHub()
	hub.SetDisconnectGracePeriod(25 * time.Millisecond)

	updates := make(chan PresenceUpdate, 4)
	hub.SetPresenceHandler(func(lobbyCode string, update PresenceUpdate) {
		if lobbyCode == "ABC123" {
			updates <- update
		}
	})

	firstConn := new(websocket.Conn)
	hub.Register("ABC123", firstConn, "player-1")
	assertUpdate(t, updates, PresenceUpdate{PlayerID: "player-1", Connected: true})

	hub.Unregister("ABC123", firstConn)
	assertUpdate(t, updates, PresenceUpdate{PlayerID: "player-1", Connected: false})

	secondConn := new(websocket.Conn)
	time.Sleep(10 * time.Millisecond)
	hub.Register("ABC123", secondConn, "player-1")
	assertUpdate(t, updates, PresenceUpdate{PlayerID: "player-1", Connected: true})

	select {
	case update := <-updates:
		t.Fatalf("unexpected extra update after reconnect: %+v", update)
	case <-time.After(35 * time.Millisecond):
	}

	presence := hub.Presence("ABC123", "player-1")
	if !presence.IsConnected() {
		t.Fatalf("expected player to be connected after reconnect")
	}
}

func assertUpdate(t *testing.T, updates <-chan PresenceUpdate, want PresenceUpdate) {
	t.Helper()

	select {
	case got := <-updates:
		if got != want {
			t.Fatalf("unexpected presence update: got %+v want %+v", got, want)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timed out waiting for presence update %+v", want)
	}
}
