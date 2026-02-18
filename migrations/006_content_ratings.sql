-- +goose Up

CREATE TABLE content_ratings (
    id SMALLINT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_content_ratings_default_true
    ON content_ratings (is_default)
    WHERE is_default = TRUE;

INSERT INTO content_ratings (id, code, label, description, is_default) VALUES
    (10, 'kids', 'Kids', 'Family-friendly with simple language and no mature themes.', FALSE),
    (20, 'work', 'Work', 'Corporate-safe content suitable for coworkers and team events.', TRUE),
    (30, 'adults', 'Adults', 'Permissive tier including edgy humor and adult references.', FALSE)
ON CONFLICT (id) DO NOTHING;

ALTER TABLE lobbies
    ADD COLUMN content_rating SMALLINT NOT NULL DEFAULT 20 REFERENCES content_ratings(id);

ALTER TABLE coordinates_prompts
    ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 30 REFERENCES content_ratings(id);

ALTER TABLE coordinates_axis_sets
    ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 30 REFERENCES content_ratings(id);

ALTER TABLE trivia_questions
    ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 30 REFERENCES content_ratings(id);

-- Existing seed pairs are intended to be usable in default Work lobbies.
UPDATE coordinates_prompts
SET min_rating = 20
WHERE created_by_label = 'cluster-seed-v1';

UPDATE coordinates_axis_sets
SET min_rating = 20
WHERE created_by_label = 'cluster-seed-v1';

-- +goose Down

ALTER TABLE trivia_questions
    DROP COLUMN IF EXISTS min_rating;

ALTER TABLE coordinates_axis_sets
    DROP COLUMN IF EXISTS min_rating;

ALTER TABLE coordinates_prompts
    DROP COLUMN IF EXISTS min_rating;

ALTER TABLE lobbies
    DROP COLUMN IF EXISTS content_rating;

DROP INDEX IF EXISTS idx_content_ratings_default_true;
DROP TABLE IF EXISTS content_ratings;
