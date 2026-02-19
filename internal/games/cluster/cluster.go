package cluster

import (
	"context"
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/games"
	clustertmpl "github.com/jgoodhcg/mindmeld/templates/cluster"
)

// ClusterGame implements the Game interface for Cluster.
type ClusterGame struct {
	queries  *db.Queries
	dbPool   *pgxpool.Pool
	eventBus events.Bus
}

// New creates a new ClusterGame.
func New(queries *db.Queries, dbPool *pgxpool.Pool, eventBus events.Bus) *ClusterGame {
	return &ClusterGame{
		queries:  queries,
		dbPool:   dbPool,
		eventBus: eventBus,
	}
}

// Info returns metadata about Cluster.
func (g *ClusterGame) Info() games.GameInfo {
	return games.GameInfo{
		Slug:        "cluster",
		Name:        "CLUSTER",
		Description: "Plot points, reveal the centroid",
		Ready:       true,
	}
}

// RenderContent builds Cluster content from lobby state.
func (g *ClusterGame) RenderContent(ctx context.Context, lobby db.Lobby, players []db.GetLobbyPlayersRow, player db.Player, isHost bool) templ.Component {
	var (
		hasRound       bool
		roundNumber    int32
		prompt         clustertmpl.PromptAxisView
		hasSubmitted   bool
		submittedCount int
		expectedCount  = len(players)
		revealed       bool
		dots           []clustertmpl.DotView
		centroidX      float64
		centroidY      float64
		standings      []clustertmpl.StandingView
		winners        []string
		outliers       []string
		discussionHint string
	)

	remainingPairs, err := g.countRemainingPromptAxisSets(ctx, g.dbPool, lobby.ID, lobby.ContentRating)
	if err != nil {
		log.Printf("[cluster] failed to count remaining prompt-axis pairs for lobby %s: %v", lobby.Code, err)
	}

	roundPoints := map[string]int{}
	activeRound, err := g.getLatestRound(ctx, g.dbPool, lobby.ID)
	if err == nil {
		hasRound = true
		roundNumber = activeRound.RoundNumber

		pair, pairErr := g.getPromptAxisSetForRound(ctx, g.dbPool, activeRound.ID)
		if pairErr != nil {
			log.Printf("[cluster] failed to get prompt-axis for round %d: %v", activeRound.RoundNumber, pairErr)
		} else {
			prompt = clustertmpl.PromptAxisView{
				PromptText: pair.PromptText,
				XMinLabel:  pair.XMinLabel,
				XMaxLabel:  pair.XMaxLabel,
				YMinLabel:  pair.YMinLabel,
				YMaxLabel:  pair.YMaxLabel,
			}
		}

		submissions, subErr := g.getRoundSubmissions(ctx, g.dbPool, activeRound.ID)
		if subErr != nil {
			log.Printf("[cluster] failed to get submissions for round %d: %v", activeRound.RoundNumber, subErr)
		} else {
			submittedCount = len(submissions)
			for _, sub := range submissions {
				if sub.PlayerID == player.ID {
					hasSubmitted = true
					break
				}
			}

			revealed = activeRound.CentroidX.Valid && activeRound.CentroidY.Valid
			if revealed {
				centroidX = activeRound.CentroidX.Float64
				centroidY = activeRound.CentroidY.Float64
				dots, roundPoints, winners, outliers = scoreRound(submissions, centroidX, centroidY, player.ID.String())
				discussionHint = discussionPrompt(prompt.PromptText, activeRound.RoundNumber)
			}
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("[cluster] failed to get latest round for lobby %s: %v", lobby.Code, err)
	}

	standings, standingsErr := g.getStandings(ctx, lobby.ID, players, roundPoints)
	if standingsErr != nil {
		log.Printf("[cluster] failed to build standings for lobby %s: %v", lobby.Code, standingsErr)
	}

	exhausted := strings.EqualFold(lobby.Phase, "finished") || (strings.EqualFold(lobby.Phase, "playing") && !hasRound && remainingPairs == 0)

	return clustertmpl.GameContent(
		lobby,
		players,
		isHost,
		minPlayersToStart,
		hasRound,
		roundNumber,
		prompt,
		submittedCount,
		expectedCount,
		hasSubmitted,
		revealed,
		dots,
		centroidX,
		centroidY,
		standings,
		winners,
		outliers,
		discussionHint,
		remainingPairs,
		exhausted,
	)
}

// RegisterRoutes registers Cluster-specific HTTP routes.
func (g *ClusterGame) RegisterRoutes(r chi.Router) {
	r.Post("/start", g.handleStartGame)
	r.Post("/submissions", g.handleSubmitCoordinate)
	r.Post("/skip", g.handleSkipPrompt)
	r.Post("/next", g.handleNextRound)
}

type coordinatesRound struct {
	ID              pgtype.UUID
	LobbyID         pgtype.UUID
	PromptAxisSetID pgtype.UUID
	RoundNumber     int32
	CentroidX       pgtype.Float8
	CentroidY       pgtype.Float8
	CreatedAt       pgtype.Timestamptz
}

type promptAxisSetRecord struct {
	ID         pgtype.UUID
	PromptText string
	XMinLabel  string
	XMaxLabel  string
	YMinLabel  string
	YMaxLabel  string
}

type submissionRecord struct {
	LobbyPlayerID pgtype.UUID
	PlayerID      pgtype.UUID
	Nickname      string
	X             float64
	Y             float64
}

type scoredSubmissionRecord struct {
	PlayerID   pgtype.UUID
	Nickname   string
	X          float64
	Y          float64
	CentroidX  float64
	CentroidY  float64
	RoundID    pgtype.UUID
	RoundScore int
}

func (g *ClusterGame) getLatestRound(ctx context.Context, q db.DBTX, lobbyID pgtype.UUID) (coordinatesRound, error) {
	const query = `
		SELECT id, lobby_id, prompt_axis_set_id, round_number, centroid_x, centroid_y, created_at
		FROM coordinates_rounds
		WHERE lobby_id = $1
		ORDER BY round_number DESC
		LIMIT 1
	`

	row := q.QueryRow(ctx, query, lobbyID)
	var round coordinatesRound
	err := row.Scan(
		&round.ID,
		&round.LobbyID,
		&round.PromptAxisSetID,
		&round.RoundNumber,
		&round.CentroidX,
		&round.CentroidY,
		&round.CreatedAt,
	)
	return round, err
}

func (g *ClusterGame) createRound(ctx context.Context, q db.DBTX, lobbyID pgtype.UUID, promptAxisSetID pgtype.UUID, roundNumber int32) (coordinatesRound, error) {
	const query = `
		INSERT INTO coordinates_rounds (lobby_id, prompt_axis_set_id, round_number)
		VALUES ($1, $2, $3)
		RETURNING id, lobby_id, prompt_axis_set_id, round_number, centroid_x, centroid_y, created_at
	`

	row := q.QueryRow(ctx, query, lobbyID, promptAxisSetID, roundNumber)
	var round coordinatesRound
	err := row.Scan(
		&round.ID,
		&round.LobbyID,
		&round.PromptAxisSetID,
		&round.RoundNumber,
		&round.CentroidX,
		&round.CentroidY,
		&round.CreatedAt,
	)
	return round, err
}

func (g *ClusterGame) getPromptAxisSetForRound(ctx context.Context, q db.DBTX, roundID pgtype.UUID) (promptAxisSetRecord, error) {
	const query = `
		SELECT cpas.id, cp.prompt_text, cas.x_min_label, cas.x_max_label, cas.y_min_label, cas.y_max_label
		FROM coordinates_rounds cr
		JOIN coordinates_prompt_axis_sets cpas ON cpas.id = cr.prompt_axis_set_id
		JOIN coordinates_prompts cp ON cp.id = cpas.prompt_id
		JOIN coordinates_axis_sets cas ON cas.id = cpas.axis_set_id
		WHERE cr.id = $1
	`

	row := q.QueryRow(ctx, query, roundID)
	var record promptAxisSetRecord
	err := row.Scan(
		&record.ID,
		&record.PromptText,
		&record.XMinLabel,
		&record.XMaxLabel,
		&record.YMinLabel,
		&record.YMaxLabel,
	)
	return record, err
}

func (g *ClusterGame) getRoundSubmissions(ctx context.Context, q db.DBTX, roundID pgtype.UUID) ([]submissionRecord, error) {
	const query = `
		SELECT cs.player_id, lp.player_id, lp.nickname, cs.x, cs.y
		FROM coordinates_submissions cs
		JOIN lobby_players lp ON lp.id = cs.player_id
		WHERE cs.round_id = $1
		ORDER BY lp.joined_at
	`

	rows, err := q.Query(ctx, query, roundID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]submissionRecord, 0)
	for rows.Next() {
		var item submissionRecord
		if scanErr := rows.Scan(&item.LobbyPlayerID, &item.PlayerID, &item.Nickname, &item.X, &item.Y); scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return items, nil
}

func (g *ClusterGame) upsertSubmission(ctx context.Context, q db.DBTX, roundID pgtype.UUID, lobbyPlayerID pgtype.UUID, x float64, y float64) error {
	const query = `
		INSERT INTO coordinates_submissions (round_id, player_id, x, y)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (round_id, player_id)
		DO UPDATE SET x = EXCLUDED.x, y = EXCLUDED.y, submitted_at = NOW()
	`

	_, err := q.Exec(ctx, query, roundID, lobbyPlayerID, x, y)
	return err
}

func (g *ClusterGame) countRoundSubmissions(ctx context.Context, q db.DBTX, roundID pgtype.UUID) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM coordinates_submissions
		WHERE round_id = $1
	`

	var count int
	err := q.QueryRow(ctx, query, roundID).Scan(&count)
	return count, err
}

func (g *ClusterGame) setRoundCentroid(ctx context.Context, q db.DBTX, roundID pgtype.UUID, centroidX float64, centroidY float64) error {
	const query = `
		UPDATE coordinates_rounds
		SET centroid_x = $2, centroid_y = $3
		WHERE id = $1
	`

	_, err := q.Exec(ctx, query, roundID, centroidX, centroidY)
	return err
}

func (g *ClusterGame) countRemainingPromptAxisSets(ctx context.Context, q db.DBTX, lobbyID pgtype.UUID, lobbyContentRating int16) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM coordinates_prompt_axis_sets cpas
		JOIN coordinates_prompts cp ON cp.id = cpas.prompt_id
		JOIN coordinates_axis_sets cas ON cas.id = cpas.axis_set_id
		WHERE cpas.is_active = TRUE
		  AND cp.is_active = TRUE
		  AND cas.is_active = TRUE
		  AND cp.min_rating <= $2
		  AND cas.min_rating <= $2
		  AND NOT EXISTS (
				SELECT 1
				FROM coordinates_rounds cr
				WHERE cr.lobby_id = $1
				  AND cr.prompt_axis_set_id = cpas.id
		  )
	`

	var count int
	err := q.QueryRow(ctx, query, lobbyID, lobbyContentRating).Scan(&count)
	return count, err
}

func (g *ClusterGame) getNextPromptAxisSet(ctx context.Context, q db.DBTX, lobbyID pgtype.UUID, lobbyContentRating int16) (promptAxisSetRecord, error) {
	const query = `
		SELECT cpas.id, cp.prompt_text, cas.x_min_label, cas.x_max_label, cas.y_min_label, cas.y_max_label
		FROM coordinates_prompt_axis_sets cpas
		JOIN coordinates_prompts cp ON cp.id = cpas.prompt_id
		JOIN coordinates_axis_sets cas ON cas.id = cpas.axis_set_id
		WHERE cpas.is_active = TRUE
		  AND cp.is_active = TRUE
		  AND cas.is_active = TRUE
		  AND cp.min_rating <= $2
		  AND cas.min_rating <= $2
		  AND NOT EXISTS (
				SELECT 1
				FROM coordinates_rounds cr
				WHERE cr.lobby_id = $1
				  AND cr.prompt_axis_set_id = cpas.id
		  )
		ORDER BY cpas.created_at, cpas.id
		LIMIT 1
	`

	row := q.QueryRow(ctx, query, lobbyID, lobbyContentRating)
	var record promptAxisSetRecord
	err := row.Scan(
		&record.ID,
		&record.PromptText,
		&record.XMinLabel,
		&record.XMaxLabel,
		&record.YMinLabel,
		&record.YMaxLabel,
	)
	return record, err
}

func (g *ClusterGame) getScoredSubmissionsForLobby(ctx context.Context, q db.DBTX, lobbyID pgtype.UUID) ([]scoredSubmissionRecord, error) {
	const query = `
		SELECT lp.player_id, lp.nickname, cs.x, cs.y, cr.centroid_x, cr.centroid_y, cr.id
		FROM coordinates_rounds cr
		JOIN coordinates_submissions cs ON cs.round_id = cr.id
		JOIN lobby_players lp ON lp.id = cs.player_id
		WHERE cr.lobby_id = $1
		  AND cr.centroid_x IS NOT NULL
		  AND cr.centroid_y IS NOT NULL
	`

	rows, err := q.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]scoredSubmissionRecord, 0)
	for rows.Next() {
		var item scoredSubmissionRecord
		if scanErr := rows.Scan(&item.PlayerID, &item.Nickname, &item.X, &item.Y, &item.CentroidX, &item.CentroidY, &item.RoundID); scanErr != nil {
			return nil, scanErr
		}
		item.RoundScore = CalculateRoundPoints(item.X, item.Y, item.CentroidX, item.CentroidY)
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return items, nil
}

func (g *ClusterGame) getStandings(ctx context.Context, lobbyID pgtype.UUID, players []db.GetLobbyPlayersRow, roundPoints map[string]int) ([]clustertmpl.StandingView, error) {
	totals := make(map[string]int, len(players))
	roundsPlayed := make(map[string]int, len(players))
	names := make(map[string]string, len(players))

	for _, p := range players {
		playerKey := p.PlayerID.String()
		totals[playerKey] = 0
		roundsPlayed[playerKey] = 0
		names[playerKey] = p.Nickname
	}

	scored, err := g.getScoredSubmissionsForLobby(ctx, g.dbPool, lobbyID)
	if err != nil {
		return nil, err
	}

	for _, row := range scored {
		key := row.PlayerID.String()
		totals[key] += row.RoundScore
		roundsPlayed[key] += 1
		if _, ok := names[key]; !ok {
			names[key] = row.Nickname
		}
	}

	standings := make([]clustertmpl.StandingView, 0, len(names))
	for key, nickname := range names {
		avg := 0.0
		if roundsPlayed[key] > 0 {
			avg = float64(totals[key]) / float64(roundsPlayed[key])
		}

		standings = append(standings, clustertmpl.StandingView{
			Nickname:          nickname,
			RoundPoints:       roundPoints[key],
			TotalPoints:       totals[key],
			AvgPointsPerRound: avg,
		})
	}

	sort.SliceStable(standings, func(i, j int) bool {
		if standings[i].TotalPoints == standings[j].TotalPoints {
			return strings.ToLower(standings[i].Nickname) < strings.ToLower(standings[j].Nickname)
		}
		return standings[i].TotalPoints > standings[j].TotalPoints
	})

	if len(standings) > 0 {
		top := standings[0].TotalPoints
		for i := range standings {
			standings[i].IsLeader = standings[i].TotalPoints == top
		}
	}

	return standings, nil
}

func scoreRound(submissions []submissionRecord, centroidX float64, centroidY float64, currentPlayerID string) ([]clustertmpl.DotView, map[string]int, []string, []string) {
	if len(submissions) == 0 {
		return nil, map[string]int{}, nil, nil
	}

	maxPoints := -1
	minPoints := 101
	roundPoints := make(map[string]int, len(submissions))
	dots := make([]clustertmpl.DotView, 0, len(submissions))

	for i, sub := range submissions {
		points := CalculateRoundPoints(sub.X, sub.Y, centroidX, centroidY)
		playerKey := sub.PlayerID.String()
		roundPoints[playerKey] = points
		if points > maxPoints {
			maxPoints = points
		}
		if points < minPoints {
			minPoints = points
		}

		dots = append(dots, clustertmpl.DotView{
			Nickname:        sub.Nickname,
			X:               sub.X,
			Y:               sub.Y,
			Points:          points,
			AnimationDelay:  i * 80,
			IsCurrentPlayer: playerKey == currentPlayerID,
		})
	}

	winners := make([]string, 0)
	outliers := make([]string, 0)
	for i := range dots {
		if dots[i].Points == maxPoints {
			dots[i].IsWinner = true
			winners = append(winners, dots[i].Nickname)
		}
		if dots[i].Points == minPoints {
			outliers = append(outliers, dots[i].Nickname)
		}
	}

	sort.Strings(winners)
	sort.Strings(outliers)
	return dots, roundPoints, winners, outliers
}

func discussionPrompt(promptText string, roundNumber int32) string {
	prompts := []string{
		`Which axis mattered most in your placement, and why?`,
		`What assumption did you make about where the group would land?`,
		`If you moved your point now, what changed in your thinking?`,
		`What did you optimize for first: the X-axis or the Y-axis?`,
	}

	if len(prompts) == 0 {
		return ""
	}

	index := int(roundNumber-1) % len(prompts)
	if index < 0 {
		index = 0
	}

	return `For "` + promptText + `": ` + prompts[index]
}
