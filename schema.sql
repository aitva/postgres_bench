CREATE TABLE pages (
    id         UUID PRIMARY KEY,
    updated_at TIMESTAMP NOT NULL,
    title      TEXT NOT NULL,
    "text"     TEXT NOT NULL
);

