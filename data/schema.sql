CREATE TABLE IF NOT EXISTS records (
    name TEXT PRIMARY KEY,
    value TEXT,
    metadata TEXT
);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);