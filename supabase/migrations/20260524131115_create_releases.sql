CREATE TYPE release_type  AS ENUM ('album', 'track');
CREATE TYPE special_label AS ENUM ('Best New Album', 'Best New Track', 'Best New Reissue');
CREATE TYPE source_name   AS ENUM ('Bandcamp', 'NPR', 'Pitchfork', 'Stereogum', 'XXL');

CREATE TABLE artists (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (name)
);

CREATE INDEX artists_deleted_at_idx ON artists (deleted_at);

CREATE TABLE releases (
    id           BIGSERIAL     PRIMARY KEY,
    name         TEXT          NOT NULL,
    publish_date TIMESTAMPTZ   NOT NULL,
    link         TEXT          NOT NULL DEFAULT '',
    source       source_name   NOT NULL,
    release_type release_type  NOT NULL,
    special      special_label,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

-- Dedup on non-empty links only (XXL reuses one URL for multiple tracks)
CREATE UNIQUE INDEX releases_link_unique    ON releases (link) WHERE link <> '';
CREATE INDEX releases_source_idx           ON releases (source);
CREATE INDEX releases_release_type_idx     ON releases (release_type);
CREATE INDEX releases_publish_date_idx     ON releases (publish_date DESC);
CREATE INDEX releases_special_idx          ON releases (special) WHERE special IS NOT NULL;
CREATE INDEX releases_deleted_at_idx       ON releases (deleted_at);

CREATE TABLE release_artists (
    release_id BIGINT NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
    artist_id  BIGINT NOT NULL REFERENCES artists(id)  ON DELETE CASCADE,
    PRIMARY KEY (release_id, artist_id)
);

CREATE INDEX release_artists_artist_idx ON release_artists (artist_id);
