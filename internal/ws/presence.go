package ws

import "time"

const defaultDisconnectGracePeriod = 25 * time.Second

// PlayerPresence is a snapshot of one player's connection state in a lobby.
type PlayerPresence struct {
	ConnectionCount int
	DisconnectedAt  time.Time
	GracePeriod     time.Duration
}

func (p PlayerPresence) IsConnected() bool {
	return p.ConnectionCount > 0
}

func (p PlayerPresence) IsDisconnected() bool {
	return p.ConnectionCount == 0 && !p.DisconnectedAt.IsZero()
}

func (p PlayerPresence) GraceExpiredAt(now time.Time) bool {
	return p.IsDisconnected() && now.Sub(p.DisconnectedAt) >= p.GracePeriod
}

func (p PlayerPresence) IsActiveAt(now time.Time) bool {
	if !p.IsDisconnected() {
		return true
	}
	return !p.GraceExpiredAt(now)
}

// PresenceUpdate is emitted whenever a player's connection state changes.
type PresenceUpdate struct {
	PlayerID     string
	Connected    bool
	GraceExpired bool
}
