package main

import (
	"log"
	"os"
	"path/filepath"
	"sloggolang/logger"
	"time"
)

func main() {

	// init file logger
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	logPath := filepath.Join(homeDir, "var/log/logapps/logging.json")

	err = os.WriteFile(logPath, nil, 0755)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger.Slog2Json(file)
	le := logger.NewLogEntry()

	ctx := le.Initialize("GET", "/api/v1/users", "localhost:8080", "body", "header", "signature", "user-id", "service-name", "service-id")

	logger.Error(ctx, 500, "Watdehel!! Error Sangat Brutal")

	// initializing thirdpartyLog
	tp := logger.NewThirdPartyLog("service-name", "1", logger.RESTYPE)

	tp.StartTraceWithContext(ctx)

	tp.FinishTrace(ctx, logger.ThirdPartyRequest{
		Url:           "https://google.com",
		Method:        "GET",
		StatusCode:    200,
		RequestBody:   "body",
		RequestHeader: "header",
		Response:      "response",
	}, time.Now())

	le.Finalize(ctx)
}
