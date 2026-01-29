CREATE TABLE IF NOT EXISTS refresh_tokens (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash bytea NOT NULL,
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS refresh_tokens_token_hash_uq
  ON refresh_tokens (token_hash);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx
  ON refresh_tokens (user_id);

CREATE INDEX IF NOT EXISTS refresh_tokens_expires_at_idx
  ON refresh_tokens (expires_at);
