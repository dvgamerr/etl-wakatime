@echo off

duckdb -no-stdin -init ./dump-data/1-wakatime-duckdb.sql
psql -U postgres -p 5433 -f ./dump-data/2-wakatime-psql.sql touno-io
psql -U postgres -p 5433 -f ./dump-data/3-wakatime-fix-duplicate.sql touno-io
rm -f wakatime-output.csv
