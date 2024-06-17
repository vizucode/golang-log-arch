package main

import (
	"sloggolang/logger"
	"time"
)

func main() {

	logger.Slog2Json()
	le := logger.NewLogEntry()

	ctx := le.Initialize("GET", "/api/v1/users", "localhost:8080", "body", "header", "signature", "user-id", "service-name", "service-id")

	logger.Error(ctx, 500, "Watdehel!! Error Ternyata")

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
