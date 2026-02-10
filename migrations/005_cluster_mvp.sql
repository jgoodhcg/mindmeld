-- +goose Up

CREATE TABLE coordinates_axis_sets (
    id UUID PRIMARY KEY,
    x_min_label TEXT NOT NULL,
    x_max_label TEXT NOT NULL,
    y_min_label TEXT NOT NULL,
    y_max_label TEXT NOT NULL,
    created_by_kind TEXT NOT NULL CHECK (created_by_kind IN ('system', 'developer', 'user')),
    created_by_player_id UUID NULL REFERENCES lobby_players(id),
    created_by_label TEXT NULL,
    authoring_mode TEXT NOT NULL CHECK (authoring_mode IN ('manual', 'ai_assisted', 'ai_generated', 'imported')),
    generator_provider TEXT NULL,
    generator_model TEXT NULL,
    generator_prompt_version TEXT NULL,
    generator_run_id TEXT NULL,
    provenance JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (
        (created_by_kind = 'user' AND created_by_player_id IS NOT NULL) OR
        (created_by_kind IN ('system', 'developer') AND created_by_player_id IS NULL)
    )
);

CREATE TABLE coordinates_prompts (
    id UUID PRIMARY KEY,
    prompt_text TEXT NOT NULL,
    created_by_kind TEXT NOT NULL CHECK (created_by_kind IN ('system', 'developer', 'user')),
    created_by_player_id UUID NULL REFERENCES lobby_players(id),
    created_by_label TEXT NULL,
    authoring_mode TEXT NOT NULL CHECK (authoring_mode IN ('manual', 'ai_assisted', 'ai_generated', 'imported')),
    generator_provider TEXT NULL,
    generator_model TEXT NULL,
    generator_prompt_version TEXT NULL,
    generator_run_id TEXT NULL,
    provenance JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (
        (created_by_kind = 'user' AND created_by_player_id IS NOT NULL) OR
        (created_by_kind IN ('system', 'developer') AND created_by_player_id IS NULL)
    )
);

CREATE TABLE coordinates_prompt_axis_sets (
    id UUID PRIMARY KEY,
    prompt_id UUID NOT NULL REFERENCES coordinates_prompts(id) ON DELETE CASCADE,
    axis_set_id UUID NOT NULL REFERENCES coordinates_axis_sets(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (prompt_id, axis_set_id)
);

CREATE TABLE coordinates_rounds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lobby_id UUID NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
    prompt_axis_set_id UUID NOT NULL REFERENCES coordinates_prompt_axis_sets(id),
    round_number INT NOT NULL,
    centroid_x DOUBLE PRECISION NULL CHECK (centroid_x >= 0 AND centroid_x <= 1),
    centroid_y DOUBLE PRECISION NULL CHECK (centroid_y >= 0 AND centroid_y <= 1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (lobby_id, round_number)
);

CREATE INDEX idx_coordinates_rounds_lobby ON coordinates_rounds(lobby_id);
CREATE INDEX idx_coordinates_rounds_prompt_axis_set ON coordinates_rounds(prompt_axis_set_id);

CREATE TABLE coordinates_submissions (
    round_id UUID NOT NULL REFERENCES coordinates_rounds(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES lobby_players(id),
    x DOUBLE PRECISION NOT NULL CHECK (x >= 0 AND x <= 1),
    y DOUBLE PRECISION NOT NULL CHECK (y >= 0 AND y <= 1),
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (round_id, player_id)
);

CREATE INDEX idx_coordinates_submissions_round ON coordinates_submissions(round_id);
CREATE INDEX idx_coordinates_submissions_player ON coordinates_submissions(player_id);

-- Seed axis sets and prompts with provenance metadata.
INSERT INTO coordinates_axis_sets (
    id,
    x_min_label,
    x_max_label,
    y_min_label,
    y_max_label,
    created_by_kind,
    created_by_label,
    authoring_mode,
    generator_provider,
    generator_model,
    generator_prompt_version,
    generator_run_id,
    provenance,
    is_active
) VALUES
(
    '11111111-1111-1111-1111-111111111101',
    'Traditional',
    'Experimental',
    'Relaxed',
    'Competitive',
    'developer',
    'cluster-seed-v1',
    'ai_generated',
    'openai',
    'gpt-5',
    'cluster_seed_v1',
    'seed-2026-02-08-a',
    '{"seeded_by":"developer","note":"cluster mvp axis set"}'::jsonb,
    TRUE
),
(
    '11111111-1111-1111-1111-111111111102',
    'Personal',
    'Universal',
    'Practical',
    'Aspirational',
    'developer',
    'cluster-seed-v1',
    'ai_generated',
    'openai',
    'gpt-5',
    'cluster_seed_v1',
    'seed-2026-02-08-b',
    '{"seeded_by":"developer","note":"cluster mvp axis set"}'::jsonb,
    TRUE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO coordinates_prompts (
    id,
    prompt_text,
    created_by_kind,
    created_by_label,
    authoring_mode,
    generator_provider,
    generator_model,
    generator_prompt_version,
    generator_run_id,
    provenance,
    is_active
) VALUES
(
    '22222222-2222-2222-2222-222222222201',
    'The best weekend activity for a mixed group',
    'developer',
    'cluster-seed-v1',
    'ai_generated',
    'openai',
    'gpt-5',
    'cluster_seed_v1',
    'seed-2026-02-08-p1',
    '{"seeded_by":"developer","topic":"social planning"}'::jsonb,
    TRUE
),
(
    '22222222-2222-2222-2222-222222222202',
    'A perfect team celebration after shipping a project',
    'developer',
    'cluster-seed-v1',
    'ai_generated',
    'openai',
    'gpt-5',
    'cluster_seed_v1',
    'seed-2026-02-08-p2',
    '{"seeded_by":"developer","topic":"work culture"}'::jsonb,
    TRUE
),
(
    '22222222-2222-2222-2222-222222222203',
    'An ideal onboarding experience for a new teammate',
    'developer',
    'cluster-seed-v1',
    'ai_generated',
    'openai',
    'gpt-5',
    'cluster_seed_v1',
    'seed-2026-02-08-p3',
    '{"seeded_by":"developer","topic":"onboarding"}'::jsonb,
    TRUE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO coordinates_prompt_axis_sets (
    id,
    prompt_id,
    axis_set_id,
    is_active
) VALUES
(
    '33333333-3333-3333-3333-333333333301',
    '22222222-2222-2222-2222-222222222201',
    '11111111-1111-1111-1111-111111111101',
    TRUE
),
(
    '33333333-3333-3333-3333-333333333302',
    '22222222-2222-2222-2222-222222222202',
    '11111111-1111-1111-1111-111111111101',
    TRUE
),
(
    '33333333-3333-3333-3333-333333333303',
    '22222222-2222-2222-2222-222222222203',
    '11111111-1111-1111-1111-111111111102',
    TRUE
)
ON CONFLICT (id) DO NOTHING;

-- +goose Down

DROP TABLE IF EXISTS coordinates_submissions;
DROP TABLE IF EXISTS coordinates_rounds;
DROP TABLE IF EXISTS coordinates_prompt_axis_sets;
DROP TABLE IF EXISTS coordinates_prompts;
DROP TABLE IF EXISTS coordinates_axis_sets;
