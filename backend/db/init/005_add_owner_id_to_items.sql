BEGIN;

ALTER TABLE items
ADD COLUMN IF NOT EXISTS owner_id BIGINT;

INSERT INTO users (id, email, password_hash, role)
VALUES (1, 'system@local', '!', 'admin')
ON CONFLICT (id) DO NOTHING;

UPDATE items
SET owner_id = 1
WHERE owner_id IS NULL;

ALTER TABLE items
ALTER COLUMN owner_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'items_owner_fk'
  ) THEN
    ALTER TABLE items
      ADD CONSTRAINT items_owner_fk
      FOREIGN KEY (owner_id) REFERENCES users(id);
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_items_owner_id ON items(owner_id);

COMMIT;
