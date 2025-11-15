-- name: CreateStrategy :one
INSERT INTO strategies (user_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetStrategyByID :one
SELECT * FROM strategies
WHERE id = $1 AND user_id = $2;

-- name: GetStrategiesByUserID :many
SELECT * FROM strategies
WHERE user_id = $1
ORDER BY name ASC;

-- name: UpdateStrategy :one
UPDATE strategies
SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $4
RETURNING *;

-- name: DeleteStrategy :exec
DELETE FROM strategies
WHERE id = $1 AND user_id = $2;
