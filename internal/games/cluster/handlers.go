package cluster

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jgoodhcg/mindmeld/internal/auth"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
)

func (g *ClusterGame) handleStartGame(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	lobby, err := g.queries.GetLobbyByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	participation, err := g.queries.GetPlayerParticipation(r.Context(), db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can start Cluster", http.StatusForbidden)
		return
	}

	players, err := g.queries.GetLobbyPlayers(r.Context(), lobby.ID)
	if err != nil {
		http.Error(w, "Failed to load lobby players", http.StatusInternalServerError)
		return
	}
	if len(players) < minPlayersToStart {
		http.Error(w, "Cluster needs at least 3 players to start", http.StatusBadRequest)
		return
	}

	if err = g.queries.UpdateLobbyPhase(r.Context(), db.UpdateLobbyPhaseParams{ID: lobby.ID, Phase: "playing"}); err != nil {
		log.Printf("[cluster] failed to set lobby %s to playing: %v", code, err)
		http.Error(w, "Failed to start Cluster", http.StatusInternalServerError)
		return
	}

	next, nextErr := g.getNextPromptAxisSet(r.Context(), g.dbPool, lobby.ID, lobby.ContentRating)
	if nextErr != nil {
		if errors.Is(nextErr, pgx.ErrNoRows) {
			if updateErr := g.queries.UpdateLobbyPhase(r.Context(), db.UpdateLobbyPhaseParams{ID: lobby.ID, Phase: "finished"}); updateErr != nil {
				log.Printf("[cluster] failed to set lobby %s as finished after empty pool: %v", code, updateErr)
			}
			g.eventBus.Publish(r.Context(), events.Event{Type: events.EventClusterExhausted, LobbyCode: code})
			http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
			return
		}

		log.Printf("[cluster] failed selecting next prompt-axis for lobby %s: %v", code, nextErr)
		http.Error(w, "Failed to start Cluster", http.StatusInternalServerError)
		return
	}

	roundNumber := int32(1)
	if latestRound, latestErr := g.getLatestRound(r.Context(), g.dbPool, lobby.ID); latestErr == nil {
		roundNumber = latestRound.RoundNumber + 1
	}

	_, createErr := g.createRound(r.Context(), g.dbPool, lobby.ID, next.ID, roundNumber)
	if createErr != nil {
		log.Printf("[cluster] failed creating first round for lobby %s: %v", code, createErr)
		http.Error(w, "Failed to create Cluster round", http.StatusInternalServerError)
		return
	}

	g.eventBus.Publish(r.Context(), events.Event{Type: events.EventClusterRoundStarted, LobbyCode: code})
	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *ClusterGame) handleSubmitCoordinate(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("x")), 64)
	if err != nil || x < 0 || x > 1 {
		http.Error(w, "x must be a number between 0 and 1", http.StatusBadRequest)
		return
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("y")), 64)
	if err != nil || y < 0 || y > 1 {
		http.Error(w, "y must be a number between 0 and 1", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	lobby, err := g.queries.GetLobbyByCode(ctx, code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}
	if !strings.EqualFold(lobby.Phase, "playing") {
		http.Error(w, "Cluster is not currently accepting submissions", http.StatusConflict)
		return
	}

	tx, err := g.dbPool.Begin(ctx)
	if err != nil {
		log.Printf("[cluster] failed to begin submission transaction: %v", err)
		http.Error(w, "Failed to submit coordinate", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	qtx := g.queries.WithTx(tx)

	participation, err := qtx.GetPlayerParticipation(ctx, db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil {
		http.Error(w, "Not in lobby", http.StatusForbidden)
		return
	}

	round, err := g.getLatestRound(ctx, tx, lobby.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "No active Cluster round", http.StatusBadRequest)
			return
		}
		log.Printf("[cluster] failed loading active round for %s: %v", code, err)
		http.Error(w, "Failed to submit coordinate", http.StatusInternalServerError)
		return
	}

	if round.CentroidX.Valid && round.CentroidY.Valid {
		http.Error(w, "This round is already revealed", http.StatusConflict)
		return
	}

	if err = g.upsertSubmission(ctx, tx, round.ID, participation.ID, x, y); err != nil {
		log.Printf("[cluster] failed upserting submission for lobby %s: %v", code, err)
		http.Error(w, "Failed to submit coordinate", http.StatusInternalServerError)
		return
	}

	players, err := qtx.GetLobbyPlayers(ctx, lobby.ID)
	if err != nil {
		log.Printf("[cluster] failed loading players for lobby %s: %v", code, err)
		http.Error(w, "Failed to submit coordinate", http.StatusInternalServerError)
		return
	}

	submissionCount, err := g.countRoundSubmissions(ctx, tx, round.ID)
	if err != nil {
		log.Printf("[cluster] failed counting submissions for lobby %s: %v", code, err)
		http.Error(w, "Failed to submit coordinate", http.StatusInternalServerError)
		return
	}

	revealed := submissionCount >= len(players) && len(players) > 0
	if revealed {
		submissions, subErr := g.getRoundSubmissions(ctx, tx, round.ID)
		if subErr != nil {
			log.Printf("[cluster] failed fetching submissions for centroid in lobby %s: %v", code, subErr)
			http.Error(w, "Failed to reveal round", http.StatusInternalServerError)
			return
		}

		points := make([]Point, 0, len(submissions))
		for _, sub := range submissions {
			points = append(points, Point{X: sub.X, Y: sub.Y})
		}
		centroidX, centroidY, ok := CalculateCentroid(points)
		if !ok {
			http.Error(w, "Failed to reveal round", http.StatusInternalServerError)
			return
		}

		if err = g.setRoundCentroid(ctx, tx, round.ID, centroidX, centroidY); err != nil {
			log.Printf("[cluster] failed storing centroid for lobby %s: %v", code, err)
			http.Error(w, "Failed to reveal round", http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("[cluster] failed committing submission transaction for lobby %s: %v", code, err)
		http.Error(w, "Failed to submit coordinate", http.StatusInternalServerError)
		return
	}

	eventType := events.EventClusterSubmissionUpdated
	var payload any = events.ClusterSubmissionUpdatedPayload{
		SubmittedCount: submissionCount,
		TotalPlayers:   len(players),
	}
	if revealed {
		eventType = events.EventClusterRoundRevealed
		payload = nil
	}
	g.eventBus.Publish(ctx, events.Event{Type: eventType, LobbyCode: code, Payload: payload})

	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}

func (g *ClusterGame) handleSkipPrompt(w http.ResponseWriter, r *http.Request) {
	g.advanceRound(w, r, false)
}

func (g *ClusterGame) handleNextRound(w http.ResponseWriter, r *http.Request) {
	g.advanceRound(w, r, true)
}

func (g *ClusterGame) advanceRound(w http.ResponseWriter, r *http.Request, requireRevealed bool) {
	code := chi.URLParam(r, "code")
	player := auth.GetPlayer(r.Context())
	ctx := r.Context()

	lobby, err := g.queries.GetLobbyByCode(ctx, code)
	if err != nil {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	participation, err := g.queries.GetPlayerParticipation(ctx, db.GetPlayerParticipationParams{
		LobbyID:  lobby.ID,
		PlayerID: player.ID,
	})
	if err != nil || !participation.IsHost {
		http.Error(w, "Only the host can continue the round", http.StatusForbidden)
		return
	}

	latestRound, err := g.getLatestRound(ctx, g.dbPool, lobby.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "No active round to continue", http.StatusBadRequest)
			return
		}
		log.Printf("[cluster] failed to get latest round for lobby %s: %v", code, err)
		http.Error(w, "Failed to continue Cluster", http.StatusInternalServerError)
		return
	}

	if requireRevealed && !(latestRound.CentroidX.Valid && latestRound.CentroidY.Valid) {
		http.Error(w, "Round must be revealed before continuing", http.StatusConflict)
		return
	}

	next, nextErr := g.getNextPromptAxisSet(ctx, g.dbPool, lobby.ID, lobby.ContentRating)
	if nextErr != nil {
		if errors.Is(nextErr, pgx.ErrNoRows) {
			if updateErr := g.queries.UpdateLobbyPhase(ctx, db.UpdateLobbyPhaseParams{ID: lobby.ID, Phase: "finished"}); updateErr != nil {
				log.Printf("[cluster] failed setting lobby %s to finished: %v", code, updateErr)
				http.Error(w, "Failed to finish Cluster", http.StatusInternalServerError)
				return
			}

			g.eventBus.Publish(ctx, events.Event{Type: events.EventClusterExhausted, LobbyCode: code})
			http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
			return
		}

		log.Printf("[cluster] failed selecting next prompt-axis for lobby %s: %v", code, nextErr)
		http.Error(w, "Failed to continue Cluster", http.StatusInternalServerError)
		return
	}

	_, createErr := g.createRound(ctx, g.dbPool, lobby.ID, next.ID, latestRound.RoundNumber+1)
	if createErr != nil {
		log.Printf("[cluster] failed creating round %d for lobby %s: %v", latestRound.RoundNumber+1, code, createErr)
		http.Error(w, "Failed to continue Cluster", http.StatusInternalServerError)
		return
	}

	g.eventBus.Publish(ctx, events.Event{Type: events.EventClusterRoundStarted, LobbyCode: code})
	http.Redirect(w, r, "/lobbies/"+code, http.StatusSeeOther)
}
