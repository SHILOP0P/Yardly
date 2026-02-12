BEGIN;

ALTER TABLE items
  ADD COLUMN IF NOT EXISTS blocked_at timestamptz NULL,
  ADD COLUMN IF NOT EXISTS block_reason text NULL;

CREATE INDEX IF NOT EXISTS items_blocked_at_idx ON items(blocked_at);

COMMIT;
