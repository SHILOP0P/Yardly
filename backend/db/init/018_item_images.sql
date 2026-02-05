
CREATE TABLE IF NOT EXISTS item_images (
  id         BIGSERIAL PRIMARY KEY,
  item_id    BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  url        TEXT   NOT NULL,
  sort_order INT    NOT NULL, -- 1..N, 1 = основная
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (item_id, sort_order)
);

CREATE INDEX IF NOT EXISTS item_images_item_id_idx
  ON item_images(item_id);
