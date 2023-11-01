SET client_encoding TO 'UTF8';

ALTER TABLE stash.wakatime_heartbeats ADD COLUMN _id SERIAL PRIMARY KEY;

DELETE FROM stash.wakatime_heartbeats a
USING stash.wakatime_heartbeats b
WHERE a._id > b._id AND a.id = b.id;

ALTER TABLE stash.wakatime_heartbeats DROP COLUMN _id;
ALTER TABLE stash.wakatime_heartbeats ADD UNIQUE (id);
