package main

import (
	"os"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
	if _, err := os.Stat(".env"); err != nil {
		return
	}
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
}

type WakaETLConfig struct {
	WakaKey string `arg:"--key,env:WAKA_APIKEY"`
	Output  string `arg:"-o,--output" default:"postgres"`
	Date    string `arg:"positional"`
	Dump    string `arg:"--dump" `
}

func main() {
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	var etl WakaETLConfig
	arg.MustParse(&etl)

	if etl.Dump != "" && etl.WakaKey == "" {
		sugar.Fatalw("Required environment variable WAKA_APIKEY")
	}

	sugar.Infof("%#v\n", etl)
}
