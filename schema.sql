CREATE TABLE pages (
    id         UUID PRIMARY KEY,
    updated_at TIMESTAMPTZ NOT NULL,
    title      TEXT NOT NULL,
    "text"     TEXT NOT NULL
);

