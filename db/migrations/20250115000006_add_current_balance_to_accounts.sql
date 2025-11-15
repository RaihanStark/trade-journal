-- migrate:up
ALTER TABLE accounts ADD COLUMN current_balance DECIMAL(20, 2) DEFAULT 0;

-- migrate:down
ALTER TABLE accounts DROP COLUMN current_balance;
