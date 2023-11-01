@echo off

curl -o heartbeats.json -X GET "https://wakatime.com/api/v1/users/current/heartbeats?date=2023-10-30" ^
 -H "Authorization: Basic d2FrYV9mMzg5ZTZkMS0yMDdlLTRjNDktOTUyNi1kMjYyM2NlN2I2ZDE=" ^
 -H "Content-Type: application/json"

duckdb -no-stdin -init ./fetch-data/1-wakatime-duckdb.sql
