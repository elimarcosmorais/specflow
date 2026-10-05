CREATE TABLE projects (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    name_key   TEXT NOT NULL UNIQUE,
    path       TEXT NOT NULL,
    path_key   TEXT NOT NULL UNIQUE,
    type       TEXT NOT NULL CHECK (type IN ('web', 'android', 'windows')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
