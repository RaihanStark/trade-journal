-- name: CreateAccount :one
INSERT INTO accounts (user_id, name, broker, account_number, account_type, currency, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, name, broker, account_number, account_type, currency, current_balance, is_active, created_at, updated_at;

-- name: GetAccountByID :one
SELECT id, user_id, name, broker, account_number, account_type, currency, current_balance, is_active, created_at, updated_at
FROM accounts
WHERE id = $1 AND user_id = $2;

-- name: GetAccountsByUserID :many
SELECT id, user_id, name, broker, account_number, account_type, currency, current_balance, is_active, created_at, updated_at
FROM accounts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $3,
    broker = $4,
    account_number = $5,
    account_type = $6,
    currency = $7,
    is_active = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, broker, account_number, account_type, currency, current_balance, is_active, created_at, updated_at;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET current_balance = COALESCE(current_balance, 0) + sqlc.arg(amount)::decimal,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, broker, account_number, account_type, currency, current_balance, is_active, created_at, updated_at;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1 AND user_id = $2;
