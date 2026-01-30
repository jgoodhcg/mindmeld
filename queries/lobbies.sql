-- name: CreateLobby :one
INSERT INTO lobbies (code, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetLobbyByCode :one
SELECT * FROM lobbies WHERE code = $1;

-- name: ListLobbies :many
SELECT * FROM lobbies ORDER BY created_at DESC;

-- name: DeleteLobby :exec
DELETE FROM lobbies WHERE id = $1;

-- name: GetLobbyStats :one
SELECT
    (SELECT COUNT(*) FROM lobbies) AS total_lobbies,
    (SELECT COUNT(*) FROM trivia_answers) AS total_answers,
    (SELECT COUNT(*) FROM trivia_rounds) AS total_rounds;
