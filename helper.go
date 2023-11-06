package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hairyhenderson/go-which"
)

func dateFormat(date string) string {
	dateLayout := "2006-01-02"
	now := time.Now().AddDate(0, 0, -1)
	if date != "" {
		pDate, err := time.Parse(dateLayout, date)
		if err == nil {
			now = pDate
		}
	}
	return now.Format(dateLayout)
}

func wakaHeartbeats(name string, date string) (string, error) {
	wakaEndpoint := "https://wakatime.com/api/v1/users/current/heartbeats"
	agent := fiber.Get(fmt.Sprintf("%s?date=%s", wakaEndpoint, date))
	agent.Add("Content-Type", "application/json").
		Add("Authorization", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(etl.WakaKey))))

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return "", errs[0]
	}
	if statusCode > 200 {
		return "", fmt.Errorf("heartbeats is status %d, %s", statusCode, body)
	}
	return writeFile(name, body)
}

func win64Goos(path string) string {
	if runtime.GOOS == "windows" && filepath.Ext(path) == "" {
		return path + ".exe"
	}
	return path
}

func writeFile(name string, body []byte) (string, error) {
	fp, err := os.CreateTemp("", name)
	if name != "" {
		fp, err = os.Create(name)
	}
	if err != nil {
		return "", err
	}
	_, err = fp.Write(body)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	return fp.Name(), nil
}

func cmd(name string, arg ...string) error {
	app := which.Which(win64Goos(name))
	cmd := exec.CommandContext(context.Background(), app, arg...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s %s\n%s\n%s", app, arg, err.Error(), output)
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		sugar.Fatal(err)
	}
}
