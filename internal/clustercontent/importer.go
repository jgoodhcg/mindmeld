package clustercontent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const idNamespace = "4f412321-1d71-42dc-a0be-f5a5f8df6f74"

type EntityPlan struct {
	DesiredCount         int
	ManagedExistingCount int
	CreateCount          int
	UpsertCount          int
	ReactivateCount      int
	DeactivateCount      int
}

type ImportPlan struct {
	Prompts  EntityPlan
	AxisSets EntityPlan
	Pairs    EntityPlan
}

func PromptUUID(slug string) uuid.UUID {
	return deterministicUUID("prompt:" + slug)
}

func AxisSetUUID(slug string) uuid.UUID {
	return deterministicUUID("axis:" + slug)
}

func PairUUID(promptSlug string, axisSlug string) uuid.UUID {
	return deterministicUUID("pair:" + promptSlug + "|" + axisSlug)
}

func Analyze(ctx context.Context, pool *pgxpool.Pool, lib Library, pairs []Pair) (ImportPlan, error) {
	promptIDs, axisIDs, pairIDs := desiredIDs(lib, pairs)

	existingPrompts, err := fetchPromptState(ctx, pool, lib.CreatedByLabel)
	if err != nil {
		return ImportPlan{}, err
	}
	existingAxisSets, err := fetchAxisState(ctx, pool, lib.CreatedByLabel)
	if err != nil {
		return ImportPlan{}, err
	}
	existingPairs, err := fetchPairState(ctx, pool, lib.CreatedByLabel)
	if err != nil {
		return ImportPlan{}, err
	}

	return ImportPlan{
		Prompts:  buildEntityPlan(promptIDs, existingPrompts),
		AxisSets: buildEntityPlan(axisIDs, existingAxisSets),
		Pairs:    buildEntityPlan(pairIDs, existingPairs),
	}, nil
}

func Import(ctx context.Context, pool *pgxpool.Pool, lib Library, pairs []Pair) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	promptIDs, axisIDs, pairIDs := desiredIDs(lib, pairs)

	for i, prompt := range lib.Prompts {
		id := promptIDs[i]
		provenance, _ := json.Marshal(map[string]string{
			"source": "cluster-content-import",
			"slug":   prompt.Slug,
		})
		if _, err := tx.Exec(ctx, `
			INSERT INTO coordinates_prompts (
				id, prompt_text, created_by_kind, created_by_label, authoring_mode,
				provenance, min_rating, is_active
			) VALUES ($1, $2, 'developer', $3, 'imported', $4::jsonb, $5, TRUE)
			ON CONFLICT (id) DO UPDATE SET
				prompt_text = EXCLUDED.prompt_text,
				created_by_kind = EXCLUDED.created_by_kind,
				created_by_label = EXCLUDED.created_by_label,
				authoring_mode = EXCLUDED.authoring_mode,
				provenance = EXCLUDED.provenance,
				min_rating = EXCLUDED.min_rating,
				is_active = TRUE
		`, id, prompt.Text, lib.CreatedByLabel, string(provenance), prompt.MinRating); err != nil {
			return fmt.Errorf("upsert prompt %s: %w", prompt.Slug, err)
		}
	}

	for i, axis := range lib.AxisSets {
		id := axisIDs[i]
		provenance, _ := json.Marshal(map[string]string{
			"source": "cluster-content-import",
			"slug":   axis.Slug,
		})
		if _, err := tx.Exec(ctx, `
			INSERT INTO coordinates_axis_sets (
				id, x_min_label, x_max_label, y_min_label, y_max_label,
				created_by_kind, created_by_label, authoring_mode, provenance, min_rating, is_active
			) VALUES ($1, $2, $3, $4, $5, 'developer', $6, 'imported', $7::jsonb, $8, TRUE)
			ON CONFLICT (id) DO UPDATE SET
				x_min_label = EXCLUDED.x_min_label,
				x_max_label = EXCLUDED.x_max_label,
				y_min_label = EXCLUDED.y_min_label,
				y_max_label = EXCLUDED.y_max_label,
				created_by_kind = EXCLUDED.created_by_kind,
				created_by_label = EXCLUDED.created_by_label,
				authoring_mode = EXCLUDED.authoring_mode,
				provenance = EXCLUDED.provenance,
				min_rating = EXCLUDED.min_rating,
				is_active = TRUE
		`, id, axis.XMinLabel, axis.XMaxLabel, axis.YMinLabel, axis.YMaxLabel, lib.CreatedByLabel, string(provenance), axis.MinRating); err != nil {
			return fmt.Errorf("upsert axis %s: %w", axis.Slug, err)
		}
	}

	for i, pair := range pairs {
		pairID := pairIDs[i]
		promptID := PromptUUID(pair.PromptSlug)
		axisID := AxisSetUUID(pair.AxisSlug)

		if _, err := tx.Exec(ctx, `
			INSERT INTO coordinates_prompt_axis_sets (id, prompt_id, axis_set_id, is_active)
			VALUES ($1, $2, $3, TRUE)
			ON CONFLICT (id) DO UPDATE SET
				prompt_id = EXCLUDED.prompt_id,
				axis_set_id = EXCLUDED.axis_set_id,
				is_active = TRUE
		`, pairID, promptID, axisID); err != nil {
			return fmt.Errorf("upsert pair %s|%s: %w", pair.PromptSlug, pair.AxisSlug, err)
		}
	}

	if err := deactivateStalePrompts(ctx, tx, lib.CreatedByLabel, promptIDs); err != nil {
		return err
	}
	if err := deactivateStaleAxisSets(ctx, tx, lib.CreatedByLabel, axisIDs); err != nil {
		return err
	}
	if err := deactivateStalePairs(ctx, tx, lib.CreatedByLabel, pairIDs); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deactivateStalePrompts(ctx context.Context, tx pgx.Tx, label string, keepIDs []uuid.UUID) error {
	if len(keepIDs) == 0 {
		_, err := tx.Exec(ctx, `UPDATE coordinates_prompts SET is_active = FALSE WHERE created_by_label = $1`, label)
		return err
	}
	_, err := tx.Exec(ctx, `
		UPDATE coordinates_prompts
		SET is_active = FALSE
		WHERE created_by_label = $1
		  AND NOT (id = ANY($2::uuid[]))
	`, label, keepIDs)
	return err
}

func deactivateStaleAxisSets(ctx context.Context, tx pgx.Tx, label string, keepIDs []uuid.UUID) error {
	if len(keepIDs) == 0 {
		_, err := tx.Exec(ctx, `UPDATE coordinates_axis_sets SET is_active = FALSE WHERE created_by_label = $1`, label)
		return err
	}
	_, err := tx.Exec(ctx, `
		UPDATE coordinates_axis_sets
		SET is_active = FALSE
		WHERE created_by_label = $1
		  AND NOT (id = ANY($2::uuid[]))
	`, label, keepIDs)
	return err
}

func deactivateStalePairs(ctx context.Context, tx pgx.Tx, label string, keepIDs []uuid.UUID) error {
	if len(keepIDs) == 0 {
		_, err := tx.Exec(ctx, `
			UPDATE coordinates_prompt_axis_sets cpas
			SET is_active = FALSE
			FROM coordinates_prompts cp, coordinates_axis_sets cas
			WHERE cpas.prompt_id = cp.id
			  AND cpas.axis_set_id = cas.id
			  AND cp.created_by_label = $1
			  AND cas.created_by_label = $1
		`, label)
		return err
	}
	_, err := tx.Exec(ctx, `
		UPDATE coordinates_prompt_axis_sets cpas
		SET is_active = FALSE
		FROM coordinates_prompts cp, coordinates_axis_sets cas
		WHERE cpas.prompt_id = cp.id
		  AND cpas.axis_set_id = cas.id
		  AND cp.created_by_label = $1
		  AND cas.created_by_label = $1
		  AND NOT (cpas.id = ANY($2::uuid[]))
	`, label, keepIDs)
	return err
}

func desiredIDs(lib Library, pairs []Pair) ([]uuid.UUID, []uuid.UUID, []uuid.UUID) {
	promptIDs := make([]uuid.UUID, 0, len(lib.Prompts))
	for _, prompt := range lib.Prompts {
		promptIDs = append(promptIDs, PromptUUID(prompt.Slug))
	}

	axisIDs := make([]uuid.UUID, 0, len(lib.AxisSets))
	for _, axis := range lib.AxisSets {
		axisIDs = append(axisIDs, AxisSetUUID(axis.Slug))
	}

	pairIDs := make([]uuid.UUID, 0, len(pairs))
	for _, pair := range pairs {
		pairIDs = append(pairIDs, PairUUID(pair.PromptSlug, pair.AxisSlug))
	}

	return promptIDs, axisIDs, pairIDs
}

func fetchPromptState(ctx context.Context, pool *pgxpool.Pool, label string) (map[uuid.UUID]bool, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, is_active
		FROM coordinates_prompts
		WHERE created_by_label = $1
	`, label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectState(rows)
}

func fetchAxisState(ctx context.Context, pool *pgxpool.Pool, label string) (map[uuid.UUID]bool, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, is_active
		FROM coordinates_axis_sets
		WHERE created_by_label = $1
	`, label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectState(rows)
}

func fetchPairState(ctx context.Context, pool *pgxpool.Pool, label string) (map[uuid.UUID]bool, error) {
	rows, err := pool.Query(ctx, `
		SELECT cpas.id, cpas.is_active
		FROM coordinates_prompt_axis_sets cpas
		JOIN coordinates_prompts cp ON cp.id = cpas.prompt_id
		JOIN coordinates_axis_sets cas ON cas.id = cpas.axis_set_id
		WHERE cp.created_by_label = $1
		  AND cas.created_by_label = $1
	`, label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectState(rows)
}

func collectState(rows pgx.Rows) (map[uuid.UUID]bool, error) {
	items := make(map[uuid.UUID]bool)
	for rows.Next() {
		var id uuid.UUID
		var isActive bool
		if err := rows.Scan(&id, &isActive); err != nil {
			return nil, err
		}
		items[id] = isActive
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func buildEntityPlan(desiredIDs []uuid.UUID, existing map[uuid.UUID]bool) EntityPlan {
	desiredSet := make(map[uuid.UUID]bool, len(desiredIDs))
	plan := EntityPlan{
		DesiredCount:         len(desiredIDs),
		ManagedExistingCount: len(existing),
	}

	for _, id := range desiredIDs {
		desiredSet[id] = true
		if isActive, ok := existing[id]; ok {
			plan.UpsertCount++
			if !isActive {
				plan.ReactivateCount++
			}
			continue
		}
		plan.CreateCount++
	}

	for id, isActive := range existing {
		if desiredSet[id] {
			continue
		}
		if isActive {
			plan.DeactivateCount++
		}
	}

	return plan
}

func deterministicUUID(value string) uuid.UUID {
	ns := uuid.MustParse(idNamespace)
	normalized := strings.TrimSpace(strings.ToLower(value))
	return uuid.NewSHA1(ns, []byte(normalized))
}
