BEGIN;

CREATE TABLE IF NOT EXISTS admin_events (
  id BIGSERIAL PRIMARY KEY,
  actor_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  entity_type TEXT NOT NULL,
  entity_id BIGINT NOT NULL,
  action TEXT NOT NULL,
  reason TEXT NULL,
  meta JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS admin_events_actor_idx ON admin_events(actor_user_id);
CREATE INDEX IF NOT EXISTS admin_events_entity_idx ON admin_events(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS admin_events_created_at_idx ON admin_events(created_at);

COMMIT;
