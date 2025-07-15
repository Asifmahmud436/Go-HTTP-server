-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, body, user_id;

-- name: GetChirpByID :one
SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1;

-- name: ListChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps ORDER BY created_at DESC;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
