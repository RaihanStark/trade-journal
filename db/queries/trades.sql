-- name: CreateTrade :one
INSERT INTO
    trades (
        user_id,
        account_id,
        date,
        time,
        pair,
        type,
        entry,
        exit,
        lots,
        pips,
        pl,
        rr,
        status,
        stop_loss,
        take_profit,
        notes,
        mistakes,
        amount
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14,
        $15,
        $16,
        $17,
        $18
    )
RETURNING
    *;

-- name: AddTradeStrategy :exec
INSERT INTO trade_strategies (trade_id, strategy_id) VALUES ($1, $2);

-- name: GetTradesByUserID :many
SELECT *
FROM trades
WHERE
    user_id = $1
ORDER BY date DESC, time DESC;

-- name: GetTradeByID :one
SELECT * FROM trades WHERE id = $1 AND user_id = $2;

-- name: GetTradeStrategies :many
SELECT s.*
FROM
    strategies s
    INNER JOIN trade_strategies ts ON s.id = ts.strategy_id
WHERE
    ts.trade_id = $1;

-- name: UpdateTrade :one
UPDATE trades
SET
    account_id = $2,
    date = $3,
    time = $4,
    pair = $5,
    type = $6,
    entry = $7,
    exit = $8,
    lots = $9,
    pips = $10,
    pl = $11,
    rr = $12,
    status = $13,
    stop_loss = $14,
    take_profit = $15,
    notes = $16,
    mistakes = $17,
    amount = $18,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND user_id = $19
RETURNING
    *;

-- name: DeleteTradeStrategies :exec
DELETE FROM trade_strategies WHERE trade_id = $1;

-- name: DeleteTrade :exec
DELETE FROM trades WHERE id = $1 AND user_id = $2;

-- name: GetTradesByAccountID :many
SELECT *
FROM trades
WHERE
    account_id = $1
    AND user_id = $2
ORDER BY date DESC, time DESC;

-- name: GetTradesByUserIDAndDateRange :many
SELECT *
FROM trades
WHERE
    user_id = $1
    AND date >= $2
    AND date <= $3
ORDER BY date DESC, time DESC;

-- name: GetTradesByAccountIDAndDateRange :many
SELECT *
FROM trades
WHERE
    account_id = $1
    AND user_id = $2
    AND date >= $3
    AND date <= $4
ORDER BY date DESC, time DESC;

-- name: UpdateTradeChartBefore :one
UPDATE trades
SET chart_before = $1, updated_at = NOW()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: UpdateTradeChartAfter :one
UPDATE trades
SET chart_after = $1, updated_at = NOW()
WHERE id = $2 AND user_id = $3
RETURNING *;