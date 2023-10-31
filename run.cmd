@echo off

duckdb -no-stdin -init wakatime-duckdb.sql
psql -U postgres -f wakatime-psql.sql touno-io
rm -f wakatime-output.csv
