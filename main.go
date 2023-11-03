package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	b64 "encoding/base64"

	"github.com/alexflint/go-arg"
	"github.com/gofiber/fiber/v2"
	"github.com/hairyhenderson/go-which"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func win64Goos(path string) string {
	if runtime.GOOS == "windows" && filepath.Ext(path) == "" {
		return path + ".exe"
	}
	return path
}

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

	body, err := wakaHeartbeats()
	if len(err) > 0 {
		for _, e := range err {
			sugar.Fatal(e)
		}
	}
	sugar.Debugf("%s", body)
	if etl.Output == "postgres" {

	} else {
		sugar.Fatalw("output not supported.")
	}
}

func wakaHeartbeats() ([]byte, []error) {
	wakaEndpoint := "https://wakatime.com/api/v1/users/current/heartbeats"
	dateLayout := "2006-01-02"
	now := time.Now().AddDate(0, 0, -1)
	if etl.Date != "" {
		etlDate, err := time.Parse(dateLayout, etl.Date)
		if err != nil {
			return nil, []error{err}
		}
		now = etlDate
	}
	agent := fiber.Get(fmt.Sprintf("%s?date=%s", wakaEndpoint, now.Format(dateLayout)))
	agent.Add("Content-Type", "application/json").
		Add("Authorization", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(etl.WakaKey))))

	statusCode, body, err := agent.Bytes()
	if len(err) > 0 || statusCode > 200 {
		return body, err
	}

	// pass status code and body received by the proxy
	return body, []error{}
}

func cmd(name string, arg ...string) error {
	var stdout bytes.Buffer
	cmd := exec.Command(name, arg...)
	// cmd.Stdin = strings.NewReader("and old falcon")

	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		sugar.Debugf("%q\n", stdout.String())
		return err
	}

	return nil
}
