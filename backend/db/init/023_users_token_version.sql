ALTER TABLE users
ADD COLUMN token_version bigint NOT NULL DEFAULT 1;
