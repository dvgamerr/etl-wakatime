-- Import data
CREATE TABLE wakatime AS SELECT UNNEST(days) as days FROM read_json_auto('wakatime-info.dvgamergmail.com.json', maximum_object_size=999999999);
-- Unnest array object
CREATE TABLE heartbeats AS SELECT days.date, UNNEST(days.heartbeats) heartbeats FROM wakatime;
-- Transfrom and export data
COPY (
  SELECT * FROM (
    SELECT
      date,
      heartbeats ->> 'id' as "id",
      heartbeats ->> 'user_agent_id' as "user_agent_id",
      heartbeats ->> 'branch' as "branch",
      heartbeats ->> 'category' as "category",
      heartbeats ->> 'type' as "type",
      CAST(heartbeats ->> 'time' AS integer) as 'time',
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
