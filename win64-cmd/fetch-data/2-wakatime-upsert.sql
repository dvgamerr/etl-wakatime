SET client_encoding TO 'UTF8';

CREATE TEMPORARY TABLE heartbeats (
  "date" DATE NOT NULL,
  "id" UUID NOT NULL,
  "user_agent_id" UUID NULL,
  "branch" VARCHAR NULL,
  "category" VARCHAR NULL,
  "type" VARCHAR NULL,
  "time" DECIMAL NULL,
  "dependencies" VARCHAR NULL,
  "entity" VARCHAR NULL,
  "language" VARCHAR NULL,
  "lineno" BIGINT NULL,
  "lines" BIGINT NULL,
  "project" VARCHAR NULL,
  "project_root_count" BIGINT NULL,
  "is_write" BOOLEAN NULL,
  "created_at" TIMESTAMP NULL,
  "cursorpos" BIGINT NULL
);

\COPY heartbeats FROM 'wakatime-output.csv' CSV HEADER;

-- DELETE FROM stash.wakatime_heartbeats WHERE id IN (SELECT id FROM heartbeats);
-- INSERT INTO stash.wakatime_heartbeats SELECT * FROM heartbeats;

INSERT INTO stash.wakatime_heartbeats
SELECT * FROM heartbeats
ON CONFLICT(id) DO NOTHING;
