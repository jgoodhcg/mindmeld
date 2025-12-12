BEGIN;

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS game_types (
  key text PRIMARY KEY,
  display_name text NOT NULL,
  config_schema jsonb NOT NULL
);

CREATE TABLE IF NOT EXISTS lobbies (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code varchar(12) UNIQUE NOT NULL,
  game_type text NOT NULL REFERENCES game_types(key),
  host_player_id uuid,
  config jsonb NOT NULL DEFAULT '{}'::jsonb,
  phase text NOT NULL CHECK (phase IN ('waiting','submitting','playing','finished')),
  created_at timestamptz NOT NULL DEFAULT now(),
  ended_at timestamptz
);

CREATE TABLE IF NOT EXISTS teams (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lobby_id uuid NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (lobby_id, name)
);

CREATE TABLE IF NOT EXISTS players (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lobby_id uuid NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
  user_id uuid,
  name text NOT NULL,
  team_id uuid REFERENCES teams(id) ON DELETE SET NULL,
  joined_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (lobby_id, name)
);

ALTER TABLE lobbies
  ADD CONSTRAINT lobbies_host_player_fk
  FOREIGN KEY (host_player_id) REFERENCES players(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS lobby_events (
  id bigserial PRIMARY KEY,
  lobby_id uuid NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
  actor_player_id uuid REFERENCES players(id) ON DELETE SET NULL,
  type text NOT NULL,
  payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trivia_questions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lobby_id uuid NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
  author_player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  question_text text NOT NULL,
  option_a text NOT NULL,
  option_b text NOT NULL,
  option_c text NOT NULL,
  option_d text NOT NULL,
  correct_option char(1) NOT NULL CHECK (correct_option IN ('A','B','C','D')),
  display_order int,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trivia_rounds (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lobby_id uuid NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
  question_id uuid NOT NULL REFERENCES trivia_questions(id) ON DELETE CASCADE,
  status text NOT NULL CHECK (status IN ('pending','active','closed')),
  started_at timestamptz,
  ended_at timestamptz
);

CREATE TABLE IF NOT EXISTS trivia_answers (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  round_id uuid NOT NULL REFERENCES trivia_rounds(id) ON DELETE CASCADE,
  player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  selected_option char(1) NOT NULL CHECK (selected_option IN ('A','B','C','D')),
  is_correct boolean NOT NULL,
  response_time_ms int,
  points_awarded int,
  answered_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (round_id, player_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_players_lobby ON players(lobby_id);
CREATE INDEX IF NOT EXISTS idx_teams_lobby ON teams(lobby_id);
CREATE INDEX IF NOT EXISTS idx_questions_lobby ON trivia_questions(lobby_id);
CREATE INDEX IF NOT EXISTS idx_questions_order ON trivia_questions(lobby_id, display_order);
CREATE INDEX IF NOT EXISTS idx_answers_round ON trivia_answers(round_id);
CREATE INDEX IF NOT EXISTS idx_lobby_events_lobby_created_at ON lobby_events(lobby_id, created_at);

COMMIT;
