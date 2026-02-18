-- name: CreateLobby :one
INSERT INTO lobbies (code, name, game_type, content_rating)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLobbyByCode :one
SELECT * FROM lobbies WHERE code = $1;

-- name: ListLobbies :many
SELECT * FROM lobbies ORDER BY created_at DESC;

-- name: DeleteLobby :exec
DELETE FROM lobbies WHERE id = $1;

-- name: UpdateLobbyContentRatingIfWaiting :execrows
UPDATE lobbies
SET content_rating = $2
WHERE id = $1
  AND phase = 'waiting';

-- name: GetLobbyStats :one
SELECT
    (SELECT COUNT(*) FROM lobbies) AS total_lobbies,
    (SELECT COUNT(*) FROM trivia_answers) AS total_answers,
    (SELECT COUNT(*) FROM trivia_rounds) AS total_rounds;
