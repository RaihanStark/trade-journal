-- migrate:up
ALTER TABLE trades ADD COLUMN chart_before TEXT;
ALTER TABLE trades ADD COLUMN chart_after TEXT;

-- migrate:down
ALTER TABLE trades DROP COLUMN chart_before;
ALTER TABLE trades DROP COLUMN chart_after;
