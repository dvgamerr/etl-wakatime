@echo off

curl --silent -o heartbeats.json -X GET "https://wakatime.com/api/v1/users/current/heartbeats?date=2023-11-02" ^
 -H "Authorization: Basic d2FrYV9mMzg5ZTZkMS0yMDdlLTRjNDktOTUyNi1kMjYyM2NlN2I2ZDE=" ^
 -H "Content-Type: application/json" >nul

duckdb -no-stdin -init ./fetch-data/1-wakatime-duckdb.sql
psql -Atx1 -U postgres -p 5433 -f "./fetch-data/2-wakatime-upsert.sql" touno-io

rm -f wakatime-output.csv heartbeats.json
