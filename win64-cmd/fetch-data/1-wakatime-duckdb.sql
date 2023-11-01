-- Import data
CREATE TABLE heartbeats AS SELECT UNNEST(data) as heartbeats FROM read_json_auto('heartbeats.json', maximum_object_size=999999999);
-- Transfrom and export data
COPY (
  SELECT * FROM (
    SELECT
      '2023-10-31' date,
      heartbeats ->> 'id' as "id",
      heartbeats ->> 'user_agent_id' as "user_agent_id",
      heartbeats ->> 'branch' as "branch",
      heartbeats ->> 'category' as "category",
      heartbeats ->> 'type' as "type",
      CAST(heartbeats ->> 'time' AS double) as 'time',
      heartbeats ->> 'dependencies' as "dependencies",
      heartbeats ->> 'entity' as "entity",
      heartbeats ->> 'language' as "language",
      heartbeats ->> 'lineno' as "lineno",
      CAST(heartbeats ->> 'lines' AS integer) as "lines",
      heartbeats ->> 'project' as "project",
      heartbeats ->> 'project_root_count' as "project_root_count",
      CAST(heartbeats ->> 'is_write' AS boolean) as "is_write",
      CAST(heartbeats ->> 'created_at' AS timestamp) as "created_at",
      CAST(heartbeats ->> 'cursorpos' AS integer) as "cursorpos"
    FROM heartbeats
  ) WHERE "category" NOT IN('browsing', 'debugging', 'designing') AND "type" NOT IN('domain')
) TO 'wakatime-output.csv' (HEADER, DELIMITER ',', ENCODING UTF8);
