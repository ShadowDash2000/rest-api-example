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

INSERT INTO song_library ("group", song, release_date, link, text) VALUES
(
 'Muse',
 'Supermassive Black Hole',
 '16.07.2006',
 'https://www.youtube.com/watch?v=Xsp3_a-PMTw',
 'Ooh baby, don''t you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS song_library
-- +goose StatementEnd
