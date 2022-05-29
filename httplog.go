package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rohanraj7316/logger"
)

func (h *HttpClient) errorLogging(ctx context.Context, method, url string, request []byte,
	start time.Time, err error) {

	l := time.Since(start).Round(time.Millisecond).String()
	lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s | %s", method,
		url, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), l)

	fields := []logger.Field{
		{
			Key:   "requestId",
			Value: ctx.Value("requestId"),
		},
		{
			Key:   "url",
			Value: url,
		},
		{
			Key:   "method",
			Value: method,
		},
		{
			Key:   "error",
			Value: err.Error(),
		},
		{
			Key:   "latency",
			Value: l,
		},
		{
			Key:   "request",
			Value: request,
		},
	}

	logger.Error(lStr, fields...)
}

func (h *HttpClient) successLogging(ctx context.Context, method, url string, status string,
	statusCode int, request, response []byte, start time.Time) {
	l := time.Since(start).Round(time.Millisecond).String()
	lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s | %s", method,
		url, statusCode, status, l)

	fields := []logger.Field{
		{
			Key:   "requestId",
			Value: ctx.Value("requestId"),
		},
		{
			Key:   "latency",
			Value: l,
		},
		{
			Key:   "url",
			Value: url,
		},
		{
			Key:   "method",
			Value: method,
		},
	}

	if h.reqResLogging {
		if h.reqResBodyLogging {
			fields = append(fields, []logger.Field{
				{
					Key:   "response",
					Value: response,
				},
				{
					Key:   "request",
					Value: request,
				},
			}...)
		}
		logger.Info(lStr, fields...)
	}
}
