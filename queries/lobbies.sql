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
    COUNT(*) AS total_lobbies,
    COUNT(DISTINCT lp.lobby_id) AS lobbies_with_players
FROM lobbies l
LEFT JOIN lobby_players lp ON l.id = lp.lobby_id;
