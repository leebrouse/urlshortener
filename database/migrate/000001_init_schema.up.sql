CREATE TABLE IF NOT EXISTS urls (
  "id" SERIAL PRIMARY KEY,
  "original_url" TEXT NOT NULL,
  "short_url" TEXT NOT NULL UNIQUE,
  "is_custom" BOOLEAN NOT NULL DEFAULT FALSE,
  "expires_at" TIMESTAMP NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_short_url ON urls (short_url);
CREATE INDEX idx_expires_at ON urls (expires_at);
