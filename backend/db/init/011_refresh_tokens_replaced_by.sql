ALTER TABLE refresh_tokens
ADD COLUMN IF NOT EXISTS replaced_by BIGINT NULL;

ALTER TABLE refresh_tokens
ADD CONSTRAINT refresh_tokens_replaced_by_fk
FOREIGN KEY (replaced_by) REFERENCES refresh_tokens(id);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_active_idx
ON refresh_tokens (user_id)
WHERE revoked_at IS NULL;

CREATE INDEX IF NOT EXISTS refresh_tokens_replaced_by_idx
ON refresh_tokens (replaced_by);
