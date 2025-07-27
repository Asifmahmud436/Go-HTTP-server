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

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetChirpByUserID :one
SELECT * FROM chirps
WHERE user_id = $1;