-- name: GetPlayerByDeviceToken :one
SELECT * FROM players WHERE device_token = $1;

-- name: CreatePlayer :one
INSERT INTO players (device_token, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: AddPlayerToLobby :one
INSERT INTO lobby_players (lobby_id, player_id, nickname, is_host)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLobbyPlayers :many
SELECT lp.*, p.device_token
FROM lobby_players lp
JOIN players p ON lp.player_id = p.id
WHERE lp.lobby_id = $1
ORDER BY lp.joined_at;

-- name: GetPlayerParticipation :one
SELECT * FROM lobby_players 
WHERE lobby_id = $1 AND player_id = $2;

-- name: UpdateLobbyPhase :exec
UPDATE lobbies SET phase = $2 WHERE id = $1;

-- name: CreateTriviaRound :one
INSERT INTO trivia_rounds (lobby_id, round_number)
VALUES ($1, $2)
RETURNING *;

-- name: GetActiveRound :one
SELECT * FROM trivia_rounds 
WHERE lobby_id = $1 
ORDER BY round_number DESC 
LIMIT 1;

-- name: UpdateRoundPhase :exec
UPDATE trivia_rounds SET phase = $2 WHERE id = $1;

-- name: UpdateRoundQuestionState :exec
UPDATE trivia_rounds 
SET current_question_id = $2, question_state = $3 
WHERE id = $1;

-- name: GetRoundState :one
SELECT tr.*, tq.question_text, tq.correct_answer, 
       tq.wrong_answer_1, tq.wrong_answer_2, tq.wrong_answer_3
FROM trivia_rounds tr
LEFT JOIN trivia_questions tq ON tr.current_question_id = tq.id
WHERE tr.id = $1;

-- name: GetAnswerStats :many
SELECT 
    selected_answer,
    COUNT(*) as count
FROM trivia_answers
WHERE question_id = $1
GROUP BY selected_answer;

-- name: CreateQuestion :one
INSERT INTO trivia_questions (round_id, author, question_text, correct_answer, wrong_answer_1, wrong_answer_2, wrong_answer_3)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetQuestionsForRound :many
SELECT * FROM trivia_questions WHERE round_id = $1 ORDER BY display_order;

-- name: UpdateQuestionOrder :exec
UPDATE trivia_questions SET display_order = $2 WHERE id = $1;

-- name: SubmitAnswer :one
INSERT INTO trivia_answers (question_id, player_id, selected_answer, is_correct)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAnswersForQuestion :many
SELECT ta.*, lp.nickname
FROM trivia_answers ta
JOIN players p ON ta.player_id = p.id
JOIN lobby_players lp ON lp.player_id = p.id AND lp.lobby_id = (
    SELECT tr.lobby_id FROM trivia_rounds tr
    JOIN trivia_questions tq ON tr.id = tq.round_id
    WHERE tq.id = ta.question_id
)
WHERE ta.question_id = $1;

-- name: CountAnswersForQuestion :one
SELECT COUNT(*) FROM trivia_answers WHERE question_id = $1;

-- name: GetLobbyScoreboard :many
SELECT 
    lp.player_id, 
    lp.nickname,
    COUNT(ta.id) FILTER (WHERE ta.is_correct) as score
FROM lobby_players lp
JOIN trivia_rounds tr ON tr.lobby_id = lp.lobby_id
JOIN trivia_questions tq ON tq.round_id = tr.id
LEFT JOIN trivia_answers ta ON ta.question_id = tq.id AND ta.player_id = lp.player_id
WHERE lp.lobby_id = $1
GROUP BY lp.player_id, lp.nickname
ORDER BY score DESC;

-- name: GetRoundScoreboard :many
SELECT 
    lp.player_id, 
    lp.nickname,
    COUNT(ta.id) FILTER (WHERE ta.is_correct) as score
FROM lobby_players lp
JOIN trivia_questions tq ON tq.round_id = $1
LEFT JOIN trivia_answers ta ON ta.question_id = tq.id AND ta.player_id = lp.player_id
WHERE lp.lobby_id = (SELECT lobby_id FROM trivia_rounds WHERE id = $1)
GROUP BY lp.player_id, lp.nickname
ORDER BY score DESC;
