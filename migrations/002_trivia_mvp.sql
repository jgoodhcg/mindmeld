-- +goose Up

-- Add game_type and phase to lobbies
ALTER TABLE lobbies ADD COLUMN game_type VARCHAR(50) NOT NULL DEFAULT 'trivia';
ALTER TABLE lobbies ADD COLUMN phase VARCHAR(20) NOT NULL DEFAULT 'waiting';

-- Player identity (device-based, auth-ready)
CREATE TABLE players (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_token VARCHAR(64) UNIQUE NOT NULL,
    user_id      UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_players_device ON players(device_token);
CREATE INDEX idx_players_user ON players(user_id) WHERE user_id IS NOT NULL;

-- Player participation in lobbies
CREATE TABLE lobby_players (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lobby_id  UUID NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    nickname  VARCHAR(50) NOT NULL,
    is_host   BOOLEAN NOT NULL DEFAULT false,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(lobby_id, player_id),
    UNIQUE(lobby_id, nickname)
);

CREATE INDEX idx_lobby_players_lobby ON lobby_players(lobby_id);
CREATE INDEX idx_lobby_players_player ON lobby_players(player_id);

-- Trivia rounds (allows multiple games per lobby)
CREATE TABLE trivia_rounds (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lobby_id     UUID NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
    round_number INTEGER NOT NULL,
    phase        VARCHAR(20) NOT NULL DEFAULT 'submitting',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(lobby_id, round_number)
);

CREATE INDEX idx_trivia_rounds_lobby ON trivia_rounds(lobby_id);

-- Questions submitted by players
CREATE TABLE trivia_questions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    round_id       UUID NOT NULL REFERENCES trivia_rounds(id) ON DELETE CASCADE,
    author         UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    question_text  TEXT NOT NULL,
    correct_answer VARCHAR(200) NOT NULL,
    wrong_answer_1 VARCHAR(200) NOT NULL,
    wrong_answer_2 VARCHAR(200) NOT NULL,
    wrong_answer_3 VARCHAR(200) NOT NULL,
    display_order  INTEGER,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_trivia_questions_round ON trivia_questions(round_id);
CREATE INDEX idx_trivia_questions_author ON trivia_questions(author);

-- Player answers to questions
CREATE TABLE trivia_answers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id     UUID NOT NULL REFERENCES trivia_questions(id) ON DELETE CASCADE,
    player_id       UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    selected_answer VARCHAR(200) NOT NULL,
    is_correct      BOOLEAN NOT NULL,
    answered_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(question_id, player_id)
);

CREATE INDEX idx_trivia_answers_question ON trivia_answers(question_id);
CREATE INDEX idx_trivia_answers_player ON trivia_answers(player_id);

-- +goose Down
DROP TABLE IF EXISTS trivia_answers;
DROP TABLE IF EXISTS trivia_questions;
DROP TABLE IF EXISTS trivia_rounds;
DROP TABLE IF EXISTS lobby_players;
DROP TABLE IF EXISTS players;
ALTER TABLE lobbies DROP COLUMN IF EXISTS phase;
ALTER TABLE lobbies DROP COLUMN IF EXISTS game_type;
