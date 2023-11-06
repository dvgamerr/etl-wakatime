package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/hairyhenderson/go-which"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		logger, _ = zap.NewDevelopment()
		err := godotenv.Load()
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		logger, _ = zap.NewProduction()
	}
	sugar = logger.Sugar()

	if !which.Found(win64Goos("duckdb")) {
		logger.Fatal("Required duckdb command 'https://duckdb.org/#quickinstall' ")
	}
	if !which.Found(win64Goos("psql")) {
		logger.Fatal("Required psql command 'https://www.postgresql.org/download/' ")
	}
}

type WakaETLConfig struct {
	WakaKey string `arg:"--key,env:WAKA_SECRET"`
	Output  string `arg:"-o,--output" default:"postgres"`
	Date    string `arg:"positional"`
	Dump    string `arg:"--dump" `
}

var etl WakaETLConfig

func main() {
	defer logger.Sync() // flushes buffer, if any

	arg.MustParse(&etl)

	if etl.Dump != "" && etl.WakaKey == "" {
		sugar.Fatalw("Required environment variable WAKA_SECRET")
	}
	currentDate := dateFormat(etl.Date)
	sugar.Infof("Query heartbeats '%s' from wakatime.com", currentDate)
	wakaFile, err := wakaHeartbeats("", currentDate)
	checkError(err)
	defer os.Remove(wakaFile)

	csvFile, err := writeFile("", []byte{})
	checkError(err)
	defer os.Remove(csvFile)

	initFile, err := writeFile("", []byte(fmt.Sprintf(`
		-- Import data
		CREATE TABLE heartbeats AS SELECT UNNEST(data) as heartbeats FROM read_json_auto('%s', maximum_object_size=999999999);
		-- Transfrom and export data
		COPY (
			SELECT * FROM (
				SELECT
					'%s' date,
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
		) TO '%s' (HEADER, DELIMITER ',', ENCODING UTF8);
	`, wakaFile, currentDate, csvFile)))
	checkError(err)
	defer os.Remove(initFile)

	sugar.Infoln("RUN:: duckdb (transfrom json to csv)")
	err = cmd("duckdb", "-no-stdin", "-init", initFile)
	checkError(err)

	psqlFile, err := writeFile("", []byte(fmt.Sprintf(`
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
		
		\COPY heartbeats FROM '%s' CSV HEADER;
		
		DELETE FROM stash.wakatime_heartbeats WHERE id IN (SELECT id FROM heartbeats);
		INSERT INTO stash.wakatime_heartbeats SELECT * FROM heartbeats;
		
		-- INSERT INTO stash.wakatime_heartbeats
		-- SELECT * FROM heartbeats
		-- ON CONFLICT(id) DO NOTHING;
	`, csvFile)))
	checkError(err)
	defer os.Remove(psqlFile)

	if etl.Output == "postgres" {
		sugar.Infoln("RUN:: psql (import to table)")
		err = cmd("psql", "-Atx1", "-f", psqlFile)
		checkError(err)
	} else {
		sugar.Fatalw("output not supported.")
	}
	sugar.Infoln("Complated export:", etl.Output)
}
