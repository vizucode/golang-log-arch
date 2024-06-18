package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

/*
Usage: SetLogger()

set APP_ENVIRONTMENT to "developement" for activate debug level
*/

func Slog2Json(file io.Writer) {
	levelDefault := slog.LevelInfo
	if strings.EqualFold("developement", os.Getenv("APP_ENVIRONTMENT")) {
		levelDefault = slog.LevelDebug
	}

	slogWithJsonHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: levelDefault,
	}).WithAttrs([]slog.Attr{
		slog.Group("metadata", slog.String("environtment", os.Getenv("APP_ENVIRONTMENT"))),
	})

	logger := slog.New(slogWithJsonHandler)
	slog.SetDefault(logger)
}

func NewLogEntry() *LogEntry {
	return &LogEntry{}
}

func (l *LogEntry) Initialize(
	reqMethod,
	reqEndpoint,
	reqHost,
	reqBody,
	reqHeader,
	signature,
	userId,
	serviceName,
	serviceId string,
) (resp context.Context) {

	var (
		tz        = tz()
		timestamp = time.Now().In(tz)

		lock = new(Locker)
	)

	l.RequestId = uuid.NewString()
	l.UserId = userId
	l.Signature = signature
	l.ServiceId = serviceId
	l.ServiceName = serviceName
	l.Timestamp = timestamp

	l.UpcomingRequest = UpcomingRequest{
		Host:           reqHost,
		Endpoint:       reqEndpoint,
		RequestMethod:  reqMethod,
		RequestHeaders: reqHeader,
		RequestBody:    reqBody,
	}

	ctx := context.WithValue(context.Background(), LogKey, lock)

	return ctx
}

func (l *LogEntry) Finalize(ctx context.Context) {
	lock, ok := extract(ctx)
	if !ok {
		log.Println(ok)
		return
	}

	l.OutgoingResponse = OutgoingRequest{
		UserId:          l.UserId,
		Signature:       l.Signature,
		Host:            l.UpcomingRequest.Host,
		Endpoint:        l.UpcomingRequest.Endpoint,
		ResponseMethod:  l.UpcomingRequest.RequestMethod,
		ResponseMessage: l.OutgoingResponse.ResponseMessage,
	}

	l.ExecTime = time.Since(l.Timestamp).Seconds()

	if val, ok := lock.LoadAndDelete(_LogMessages); ok && val != nil {
		l.LogMessage = val.([]LogMessage)
	}

	if val, ok := lock.LoadAndDelete(_ResponseMessage); ok && val != nil {
		l.OutgoingResponse.ResponseMessage = val.(string)
	}

	if val, ok := lock.LoadAndDelete(_StatusCode); ok && val != nil {
		l.OutgoingResponse.StatusCode = val.(int)
	}

	if val, ok := lock.LoadAndDelete(_ThirdParties); ok && val != nil {
		l.ThirdPartyService = append(l.ThirdPartyService, val.([]ThirdpartyLog)...)
	}

	// check status code
	switch l.OutgoingResponse.StatusCode {
	case 400:
		slog.WarnContext(ctx, "data", slog.Any("data", l.serialize()))
	case 500:
		slog.ErrorContext(ctx, "data", slog.Any("data", l.serialize()))
	default:
		slog.InfoContext(ctx, "data", slog.Any("data", l.serialize()))
	}
}

func (l *LogEntry) serialize() (resp map[string]interface{}) {
	// serialization
	jByte, err := json.Marshal(l)
	if err != nil {
		return
	}

	var unmarshaler map[string]interface{}
	err = json.Unmarshal(jByte, &unmarshaler)
	if err != nil {
		return
	}

	return unmarshaler
}
