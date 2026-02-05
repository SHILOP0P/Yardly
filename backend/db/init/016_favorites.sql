CREATE TABLE IF NOT EXISTS favorites (
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  item_id BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, item_id)
);

-- ускоряет “кто добавил этот item”
CREATE INDEX IF NOT EXISTS favorites_item_id_idx ON favorites(item_id);
