package logger

import (
	"sync"
	"time"
)

type Locker struct {
	data sync.Map
}

type (
	// Key of context
	Key int
	// Flags is key for store context
	Flags string
	// ServiceType is type for data logging
	ServiceType string

	APIType string
)

const (
	// logKey is key context for request http rest api
	LogKey = Key(31)

	// Flags for key of struct
	_LogMessages     Flags = "LogMessages"
	_StatusCode      Flags = "StatusCode"
	_ResponseMessage Flags = "ResponseMessage"
	_ThirdParties    Flags = "ThirdParties"

	// API TYPE
	GRPCTYPE APIType = "grpc"
	RESTYPE  APIType = "rest"

	// type level
	levelDebug string = "debug"
	levelError string = "error"
	levelInfo  string = "info"
	levelWarn  string = "warning"
)

type UpcomingRequest struct {
	Host           string `json:"host"`
	Endpoint       string `json:"endpoint"`
	RequestMethod  string `json:"request_method"`
	RequestHeaders string `json:"request_headers"`
	RequestBody    string `json:"request_body"`
}

type OutgoingRequest struct {
	UserId          string `json:"user_id"`
	Signature       string `json:"signature"`
	Host            string `json:"host"`
	Endpoint        string `json:"endpoint"`
	ResponseMethod  string `json:"response_method"`
	StatusCode      int    `json:"status_code"`
	ResponseMessage string `json:"response_message"` // this's will encoded with base64
}

type ThirdPartyRequest struct {
	APIType       APIType `json:"api_type"`
	Url           string  `json:"url"`
	Method        string  `json:"method"`
	RequestHeader string  `json:"request_header"`
	StatusCode    int     `json:"status_code"`
	RequestBody   string  `json:"request_body"`
	Response      string  `json:"response"`
}

type ThirdpartyLog struct {
	ApiType       APIType   `json:"api_type"`
	ServiceName   string    `json:"service_name"` // required
	ServiceId     string    `json:"service_id"`
	RequestTime   time.Time `json:"request_time"`
	ExecutedTime  float64   `json:"executed_time"`
	Url           string    `json:"url"`
	Method        string    `json:"method"`
	RequestHeader string    `json:"request_header"`
	StatusCode    int       `json:"status_code"`
	RequestBody   string    `json:"request_body"`
	Response      string    `json:"response"`
}

type LogMessage struct {
	File    string `json:"file"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type LogEntry struct {
	ServiceId   string    `json:"service_id"`
	RequestId   string    `json:"request_id"`
	ServiceName string    `json:"service_name"`
	UserId      string    `json:"user_id"`
	Signature   string    `json:"signature"`
	Timestamp   time.Time `json:"timestamp"`
	ExecTime    float64   `json:"exec_time"`

	LogMessage        []LogMessage    `json:"log_message"`
	UpcomingRequest   UpcomingRequest `json:"upcoming_request"`
	OutgoingResponse  OutgoingRequest `json:"outgoing_response"`
	ThirdPartyService []ThirdpartyLog `json:"third_party_service"`
}
