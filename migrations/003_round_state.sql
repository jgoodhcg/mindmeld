-- +goose Up
ALTER TABLE trivia_rounds ADD COLUMN current_question_id UUID REFERENCES trivia_questions(id);
ALTER TABLE trivia_rounds ADD COLUMN question_state VARCHAR(20) NOT NULL DEFAULT 'idle'; -- idle, answering, revealed

-- +goose Down
ALTER TABLE trivia_rounds DROP COLUMN IF EXISTS question_state;
ALTER TABLE trivia_rounds DROP COLUMN IF EXISTS current_question_id;
