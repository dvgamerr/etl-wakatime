package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alexflint/go-arg"
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

func main() {
	defer logger.Sync() // flushes buffer, if any

	var etl WakaETLConfig
	arg.MustParse(&etl)

	if etl.Dump != "" && etl.WakaKey == "" {
		sugar.Fatalw("Required environment variable WAKA_SECRET")
	}

	if etl.Output == "postgres" {

	} else {
		sugar.Fatalw("output not supported.")
	}
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
