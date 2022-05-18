package httpclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rohanraj7316/logger"
)

type HttpClient struct {
	client            *http.Client
	reqResLogging     bool
	reqResBodyLogging bool
}

func New(config ...Config) (*HttpClient, error) {
	cfg := configDefault(config...)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// setting up proxy
	if cfg.UseProxy {
		pProxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, err
		}

		transport.Proxy = http.ProxyURL(pProxyURL)
	}

	err := logger.Configure()
	if err != nil {
		return nil, err
	}

	return &HttpClient{
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		reqResLogging:     cfg.LogReqResEnable,
		reqResBodyLogging: cfg.LogReqResBodyEnable,
	}, nil
}

func (h *HttpClient) successLogging(method, url, status string, statusCode int, response io.ReadCloser,
	start time.Time, fields ...logger.Field) {
	l := time.Since(start).Round(time.Millisecond).String()
	if h.reqResLogging {
		lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s | %s", method, url, statusCode, status, l)
		if h.reqResBodyLogging {
			var responseBody interface{}
			if response != nil {
				err := json.NewDecoder(response).Decode(&responseBody)
				if err != nil {
					logger.Error(err.Error())
				}
			}

			fields = append(fields, []logger.Field{
				{
					Key:   "statusCode",
					Value: statusCode,
				},
				{
					Key:   "response",
					Value: responseBody,
				},
				{
					Key:   "latency",
					Value: l,
				},
			}...)
			logger.Info(lStr, fields...)
		} else {
			logger.Info(lStr)
		}
	}
}

func (h *HttpClient) errorLogging(method, url string, start time.Time, err error, fields ...logger.Field) {
	l := time.Since(start).Round(time.Millisecond).String()
	lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s | %s", method,
		url, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), l)
	fields = append(fields, []logger.Field{
		{
			Key:   "error",
			Value: err.Error(),
		},
		{
			Key:   "latency",
			Value: l,
		},
	}...)
	logger.Error(lStr, fields...)
}

// Request responsible for sending http request
// by using the Option set at the time of initialization.
func (h *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string,
	request io.Reader) (*http.Response, error) {
	start := time.Now()
	lBody := []logger.Field{
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
	}

	if h.reqResBodyLogging {
		if request != nil {
			bRequest, err := io.ReadAll(request)
			if err != nil {
				logger.Error(err.Error())
			} else {
				lBody = append(lBody, logger.Field{
					Key:   "request",
					Value: bRequest,
				})
			}
		}

		if headers != nil {
			lBody = append(lBody, logger.Field{
				Key:   "headers",
				Value: headers,
			})
		}
	}

	rBody, err := http.NewRequestWithContext(ctx, method, url, request)
	if err != nil {
		h.errorLogging(method, url, start, err, lBody...)
		return nil, err
	}

	for key, val := range headers {
		rBody.Header[key] = []string{val}
	}

	response, err := h.client.Do(rBody)
	if err != nil {
		h.errorLogging(method, url, start, err, lBody...)
		return nil, err
	}

	h.successLogging(method, url, response.Status, response.StatusCode, response.Body, start, lBody...)

	return response, nil
}

// TODO: SOAP REQUEST
