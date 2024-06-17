package logger

import (
	"context"
	"time"
)

func NewThirdPartyLog(serviceName string, serviceId string, ApiType APIType) *ThirdpartyLog {
	return &ThirdpartyLog{
		ServiceName: serviceName,
		ServiceId:   serviceId,
		ApiType:     ApiType,
	}
}

func (tp *ThirdpartyLog) StartTraceWithContext(ctx context.Context) {
	tp.RequestTime = time.Now().In(tz())
}

func (tp *ThirdpartyLog) FinishTrace(ctx context.Context, payload ThirdPartyRequest, timeStart time.Time) {
	tp.Url = payload.Url
	tp.Method = payload.Method
	tp.StatusCode = payload.StatusCode
	tp.RequestBody = payload.RequestBody
	tp.RequestHeader = payload.RequestHeader
	tp.Response = payload.Response
	tp.ExecutedTime = time.Since(timeStart.In(tz())).Seconds()

	tp.store(ctx)
}

func (tp *ThirdpartyLog) store(ctx context.Context) {

	var (
		data []ThirdpartyLog
	)

	// extract sync.Map from context
	val, ok := extract(ctx)
	if !ok {
		return
	}

	// load and delete third party on sync.Map
	tmp, ok := val.LoadAndDelete(_ThirdParties)
	if ok {
		data = tmp.([]ThirdpartyLog)
	}

	// append current third party to sync.Map
	data = append(data, *tp)
	val.Set(_ThirdParties, data)
}
