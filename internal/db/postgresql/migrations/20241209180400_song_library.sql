-- +goose Up
-- +goose StatementBegin
CREATE TABLE song_library
(
    id           SERIAL PRIMARY KEY,
    "group"      VARCHAR(255) NOT NULL,
    song         VARCHAR(255) NOT NULL,
    release_date VARCHAR(10),
    link         TEXT,
    text         TEXT
);

ALTER TABLE song_library
    ADD CONSTRAINT unique_group_song
        UNIQUE ("group", song);

CREATE INDEX idx_song_library_group ON song_library ("group");
CREATE INDEX idx_song_library_song ON song_library (song);
CREATE INDEX idx_song_library_text ON song_library USING GIN (to_tsvector('simple', text));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS song_library
-- +goose StatementEnd
