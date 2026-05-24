ALTER TABLE artists DROP CONSTRAINT artists_name_key;
CREATE UNIQUE INDEX artists_name_lower_unique ON artists (LOWER(name));
