ALTER TABLE items
  ADD COLUMN description text NOT NULL DEFAULT '',
  ADD COLUMN price bigint NOT NULL DEFAULT 0,
  ADD COLUMN deposit bigint NOT NULL DEFAULT 0,
  ADD COLUMN location text NOT NULL DEFAULT '',
  ADD COLUMN category text NOT NULL DEFAULT '';
