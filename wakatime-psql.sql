SET client_encoding TO 'UTF8';

DROP TABLE IF EXISTS heartbeats;
CREATE TABLE heartbeats (
  "date" DATE NOT NULL,
  "id" UUID NULL,
  "user_agent_id" UUID NULL,
  "branch" VARCHAR NULL,
  "category" VARCHAR NULL,
  "type" VARCHAR NULL,
  "time" BIGINT NULL,
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
