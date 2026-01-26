CREATE TABLE IF NOT EXISTS items (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  status TEXT NOT NULL,
  mode TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings (
  id BIGSERIAL PRIMARY KEY,
  item_id BIGINT NOT NULL REFERENCES items(id),
  requester_id BIGINT NOT NULL,
  owner_id BIGINT NOT NULL,

  type TEXT NOT NULL,
  status TEXT NOT NULL,

  start_at TIMESTAMPTZ NULL,
  end_at   TIMESTAMPTZ NULL,

  handover_deadline TIMESTAMPTZ NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_bookings_item_id ON bookings(item_id);
CREATE INDEX IF NOT EXISTS idx_bookings_item_status ON bookings(item_id, status);
