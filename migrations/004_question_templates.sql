-- +goose Up

-- Track which templates have been used in a lobby (across all rounds)
CREATE TABLE used_question_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lobby_id    UUID NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
    template_id VARCHAR(50) NOT NULL,
    used_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(lobby_id, template_id)
);

CREATE INDEX idx_used_templates_lobby ON used_question_templates(lobby_id);

-- +goose Down
DROP TABLE IF EXISTS used_question_templates;
