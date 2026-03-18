package lobbyview

import "github.com/jgoodhcg/mindmeld/internal/db"

// HostTransferOption is an eligible transfer target for the current host.
type HostTransferOption struct {
	PlayerID string
	Nickname string
}

func BuildHostTransferOptions(
	players []db.GetLobbyPlayersRow,
	currentPlayerID string,
	isConnected func(playerID string) bool,
) []HostTransferOption {
	options := make([]HostTransferOption, 0, len(players))
	for _, player := range players {
		playerID := player.PlayerID.String()
		if playerID == currentPlayerID {
			continue
		}
		if !isConnected(playerID) {
			continue
		}
		options = append(options, HostTransferOption{
			PlayerID: playerID,
			Nickname: player.Nickname,
		})
	}
	return options
}
