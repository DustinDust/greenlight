CREATE INDEX IF NOT EXISTS movies_title_index ON movies USING GIN(to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_geners_index ON movies USING GIN(genres);