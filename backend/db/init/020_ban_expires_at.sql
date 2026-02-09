ALTER TABLE users
  ADD COLUMN IF NOT EXISTS ban_expires_at timestamptz NULL;

CREATE INDEX IF NOT EXISTS users_ban_expires_at_idx ON users(ban_expires_at);
