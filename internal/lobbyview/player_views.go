package lobbyview

import (
	"github.com/jgoodhcg/mindmeld/internal/db"
)

// Presence is the UI-relevant presence state for one player.
type Presence struct {
	Disconnected bool
	GraceExpired bool
}

// Player is the UI-facing representation of one lobby participant.
type Player struct {
	Nickname     string
	IsHost       bool
	Disconnected bool
	GraceExpired bool
}

func Build(players []db.GetLobbyPlayersRow, presence map[string]Presence) []Player {
	views := make([]Player, 0, len(players))
	for _, player := range players {
		state := presence[player.PlayerID.String()]
		views = append(views, Player{
			Nickname:     player.Nickname,
			IsHost:       player.IsHost,
			Disconnected: state.Disconnected,
			GraceExpired: state.GraceExpired,
		})
	}
	return views
}
