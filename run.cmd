@echo off

duckdb -no-stdin -init ./dump-data/wakatime-duckdb.sql
psql -U postgres -p 5433 -f ./dump-data/wakatime-psql.sql touno-io
rm -f wakatime-output.csv
