package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/contentrating"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/lobbyview"
	"github.com/jgoodhcg/mindmeld/templates"
)

func (s *Server) handleCreateLobby(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	lobbyName := strings.TrimSpace(r.FormValue("name"))
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	gameType := strings.TrimSpace(r.FormValue("game_type"))
	contentRating, err := contentRatingFromForm(r)
	if err != nil {
		http.Error(w, "Invalid content rating", http.StatusBadRequest)
		return
	}

	if lobbyName == "" || nickname == "" {
		http.Error(w, "Lobby name and nickname are required", http.StatusBadRequest)
		return
	}

	if gameType == "" {
		gameType = "trivia"
	}

	player := auth.GetPlayer(r.Context())
	code := generateCode()

	// Create Lobby
	lobby, err := s.queries.CreateLobby(r.Context(), db.CreateLobbyParams{
		Code:          code,
		Name:          lobbyName,
		GameType:      gameType,
		ContentRating: contentRating,
	})
	if err != nil {
		log.Printf("Error creating lobby: %v", err)
		http.Error(w, "Failed to create lobby", http.StatusInternalServerError)
		return
	}

	// Add Host Player
	_, err = s.queries.AddPlayerToLobby(r.Context(), db.AddPlayerToLobbyParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
		Nickname: nickname,
		IsHost:   true,
	})
	if err != nil {
		log.Printf("Error adding host to lobby: %v", err)
		http.Error(w, "Failed to join lobby", http.StatusInternalServerError)
		return
	}

	// Publish event for real-time updates (host joining counts as a player join)
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventPlayerJoined,
		LobbyCode: code,
		Payload: events.PlayerJoinedPayload{
			PlayerID: player.ID.String(),
			Nickname: nickname,
		},
	})

	http.Redirect(w, r, "/lobbies/"+lobby.Code, http.StatusSeeOther)
}

func (s *Server) handleLobbyRoom(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	// Check if player is already participating
	participation, err := s.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})

	// If checking failed (not found), show Join screen
	if err != nil {
		templates.JoinLobby(lobby).Render(r.Context(), w)
		return
	}

	// Player is in the lobby, show the room
	players, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		log.Printf("Error fetching lobby players: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Look up game by lobby's game_type and render content
	game, ok := s.games.Get(lobby.GameType)
	if !ok {
		log.Printf("Unknown game type: %s", lobby.GameType)
		http.Error(w, "Unknown game type", http.StatusInternalServerError)
		return
	}

	gameContent := game.RenderContent(r.Context(), lobby, players, player, participation.IsHost)
	rawPresence := s.hub.Snapshot(code)
	now := time.Now()
	presence := make(map[string]lobbyview.Presence, len(rawPresence))
	for playerID, state := range rawPresence {
		presence[playerID] = lobbyview.Presence{
			Disconnected: state.IsDisconnected(),
			GraceExpired: state.GraceExpiredAt(now),
		}
	}
	playerViews := lobbyview.Build(players, presence)
	templates.LobbyRoom(lobby, playerViews, participation.IsHost, player.ID.String(), gameContent).Render(r.Context(), w)
}

func (s *Server) handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	nickname := strings.TrimSpace(r.FormValue("nickname"))

	if nickname == "" {
		http.Error(w, "Nickname is required", http.StatusBadRequest)
		return
	}

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	player := auth.GetPlayer(r.Context())

	// Add player to lobby
	_, err = s.queries.AddPlayerToLobby(r.Context(), db.AddPlayerToLobbyParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
		Nickname: nickname,
		IsHost:   false,
	})
	if err != nil {
		// If duplicate key error, maybe they refreshed or clicked twice?
		// For now, log and fail, or we could redirect if they are already in.
		log.Printf("Error joining lobby: %v", err)
		http.Error(w, "Failed to join lobby (Name taken?)", http.StatusInternalServerError)
		return
	}

	// Publish event for real-time updates
	s.eventBus.Publish(r.Context(), events.Event{
		Type:      events.EventPlayerJoined,
		LobbyCode: code,
		Payload: events.PlayerJoinedPayload{
			PlayerID: player.ID.String(),
			Nickname: nickname,
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (s *Server) handleTransferHost(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())
	targetPlayerID := strings.TrimSpace(r.FormValue("target_player_id"))
	if targetPlayerID == "" {
		http.Error(w, "Target player is required", http.StatusBadRequest)
		return
	}

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	participation, err := s.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can transfer host", http.StatusForbidden)
		return
	}

	allowed, err := s.manualHostTransferAllowed(r.Context(), lobby)
	if err != nil {
		log.Printf("Error checking manual host transfer state for lobby %s: %v", code, err)
		http.Error(w, "Failed to transfer host", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Host transfer is only available before the game starts or between rounds", http.StatusConflict)
		return
	}

	players, err := s.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "Failed to load lobby players", http.StatusInternalServerError)
		return
	}

	var target *db.GetLobbyPlayersRow
	for i := range players {
		if players[i].PlayerID.String() == targetPlayerID {
			target = &players[i]
			break
		}
	}
	if target == nil {
		http.Error(w, "Selected player is not in this lobby", http.StatusBadRequest)
		return
	}
	if target.PlayerID == player.ID {
		http.Error(w, "Transfer target must be another player", http.StatusBadRequest)
		return
	}
	if !s.hub.Presence(code, target.PlayerID.String()).IsConnected() {
		http.Error(w, "Transfer target must be connected", http.StatusConflict)
		return
	}

	ctx := r.Context()
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		http.Error(w, "Failed to transfer host", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `UPDATE lobby_players SET is_host = false WHERE lobby_id = $1`, lobby.ID); err != nil {
		http.Error(w, "Failed to transfer host", http.StatusInternalServerError)
		return
	}
	if _, err = tx.Exec(ctx, `UPDATE lobby_players SET is_host = true WHERE lobby_id = $1 AND player_id = $2`, lobby.ID, target.PlayerID); err != nil {
		http.Error(w, "Failed to transfer host", http.StatusInternalServerError)
		return
	}
	if err = tx.Commit(ctx); err != nil {
		http.Error(w, "Failed to transfer host", http.StatusInternalServerError)
		return
	}

	s.eventBus.Publish(ctx, events.Event{
		Type:      events.EventHostTransferred,
		LobbyCode: code,
		Payload: events.HostTransferredPayload{
			FromPlayerID: player.ID.String(),
			ToPlayerID:   target.PlayerID.String(),
		},
	})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (s *Server) handleUpdateLobbyContentRating(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := s.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	participation, err := s.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can change audience", http.StatusForbidden)
		return
	}

	contentRating, err := contentRatingFromForm(r)
	if err != nil {
		http.Error(w, "Invalid content rating", http.StatusBadRequest)
		return
	}

	rowsUpdated, err := s.queries.UpdateLobbyContentRatingIfWaiting(r.Context(), db.UpdateLobbyContentRatingIfWaitingParams{
		ID:            lobby.ID,
		ContentRating: contentRating,
	})
	if err != nil {
		log.Printf("Error updating content rating: %v", err)
		http.Error(w, "Failed to update audience", http.StatusInternalServerError)
		return
	}
	if rowsUpdated == 0 {
		http.Error(w, "Audience can only be changed before game start", http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func contentRatingFromForm(r *http.Request) (int16, error) {
	if raw := strings.TrimSpace(r.FormValue("content_rating")); raw != "" {
		return contentrating.ParseID(raw)
	}
	return contentrating.FromPoliteMode(r.FormValue("polite_mode") != ""), nil
}

func (s *Server) manualHostTransferAllowed(ctx context.Context, lobby db.Lobby) (bool, error) {
	switch lobby.Phase {
	case "waiting":
		return true, nil
	case "playing":
		switch lobby.GameType {
		case "trivia":
			round, err := s.queries.GetActiveRound(ctx, lobby.ID)
			if err != nil {
				return false, err
			}
			if round.Phase == "submitting" || round.Phase == "finished" {
				return true, nil
			}
			return round.Phase == "playing" && round.QuestionState == "revealed", nil
		case "cluster":
			var revealed bool
			err := s.dbPool.QueryRow(ctx, `
				SELECT centroid_x IS NOT NULL AND centroid_y IS NOT NULL
				FROM coordinates_rounds
				WHERE lobby_id = $1
				ORDER BY round_number DESC
				LIMIT 1
			`, lobby.ID).Scan(&revealed)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return false, nil
				}
				return false, err
			}
			return revealed, nil
		default:
			return false, nil
		}
	default:
		return false, nil
	}
}
