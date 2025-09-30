-- +goose Up

CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS url_stats (
    id SERIAL PRIMARY KEY,
    url_id INTEGER NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    user_id TEXT,
    user_agent TEXT,
    ip_address INET,
    referer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_urls_alias ON urls(alias);
CREATE INDEX IF NOT EXISTS idx_url_stats_url_id ON url_stats(url_id);
CREATE INDEX IF NOT EXISTS idx_url_stats_created_at ON url_stats(created_at);
CREATE INDEX IF NOT EXISTS idx_url_stats_user_agent ON url_stats(user_agent);

-- +goose Down
DROP INDEX IF EXISTS idx_url_stats_user_agent;
DROP INDEX IF EXISTS idx_url_stats_created_at;
DROP INDEX IF EXISTS idx_url_stats_url_id;
DROP INDEX IF EXISTS idx_urls_alias;
DROP TABLE IF EXISTS url_stats;
DROP TABLE IF EXISTS urls;