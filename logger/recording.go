package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

func Error(ctx context.Context, statusCode int, msg ...string) {

	var (
		file     string
		messages []LogMessage
	)
	// for get filename and line when developer called this method
	_, fileName, line, _ := runtime.Caller(1)
	files := strings.Split(fileName, "/")
	if len(files) > 3 {
		file = fmt.Sprintf("%s:%d", strings.Join(files[len(files)-2:], "/"), line)
	} else {
		file = fmt.Sprintf("%s:%d", strings.Join(files, "/"), line)
	}

	val, ok := extract(ctx)
	if !ok {
		return
	}

	tmp, ok := val.LoadAndDelete(_LogMessages)
	if ok {
		messages = tmp.([]LogMessage)
	}

	messages = append(messages, LogMessage{
		File:    file,
		Level:   levelError,
		Message: fmt.Sprint(msg),
	})

	val.Set(_LogMessages, messages)
	val.Set(_StatusCode, statusCode)
	val.Set(_ResponseMessage, fmt.Sprint(msg))
}
