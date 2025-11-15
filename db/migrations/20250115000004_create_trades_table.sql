-- migrate:up
CREATE TYPE trade_type AS ENUM ('BUY', 'SELL', 'DEPOSIT', 'WITHDRAW');
CREATE TYPE trade_status AS ENUM ('open', 'closed');

CREATE TABLE IF NOT EXISTS trades (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    time TIME NOT NULL,
    pair VARCHAR(20),
    type trade_type NOT NULL,
    entry DECIMAL(20, 8),
    exit DECIMAL(20, 8),
    lots DECIMAL(10, 2),
    pips DECIMAL(10, 2),
    pl DECIMAL(20, 2),
    rr DECIMAL(10, 2),
    status trade_status NOT NULL DEFAULT 'open',
    stop_loss DECIMAL(20, 8),
    take_profit DECIMAL(20, 8),
    notes TEXT,
    mistakes TEXT,
    amount DECIMAL(20, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS trade_strategies (
    trade_id INTEGER NOT NULL REFERENCES trades(id) ON DELETE CASCADE,
    strategy_id INTEGER NOT NULL REFERENCES strategies(id) ON DELETE CASCADE,
    PRIMARY KEY (trade_id, strategy_id)
);

-- migrate:down
DROP TABLE IF EXISTS trade_strategies;
DROP TABLE IF EXISTS trades;
DROP TYPE IF EXISTS trade_status;
DROP TYPE IF EXISTS trade_type;
