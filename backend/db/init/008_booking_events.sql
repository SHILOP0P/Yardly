CREATE TABLE IF NOT EXISTS booking_events (
  id bigserial PRIMARY KEY,
  booking_id bigint NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
  actor_user_id bigint NULL REFERENCES users(id),
  action text NOT NULL,
  from_status text NULL,
  to_status text NULL,
  meta jsonb NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS booking_events_booking_id_created_at_idx
  ON booking_events (booking_id, created_at DESC, id DESC);
