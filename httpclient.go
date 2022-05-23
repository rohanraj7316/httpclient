package httpclient

import (
	"bytes"
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

func (h *HttpClient) errStr(msg string, err error) string {
	return fmt.Sprintf("[HttpClient] %s: %s", msg, err.Error())
}

func (h *HttpClient) successLogging(method, url string, request *bytes.Buffer, status string, statusCode int, response io.ReadCloser,
	start time.Time, fields ...logger.Field) {
	l := time.Since(start).Round(time.Millisecond).String()
	if h.reqResLogging {
		lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s | %s", method, url, statusCode, status, l)
		fields := append(fields, []logger.Field{
			{
				Key:   "statusCode",
				Value: statusCode,
			},
			{
				Key:   "latency",
				Value: l,
			},
		}...)

		if h.reqResBodyLogging {
			var resBody interface{}
			if response != nil {
				err := json.NewDecoder(response).Decode(&resBody)
				if err != nil {
					fmt.Println(resBody)
					logger.Error(h.errStr("failed to parse response", err))
				}
			}

			var reqByte []byte
			if request != nil {
				rByte, err := io.ReadAll(request)
				if err != nil {
					logger.Error(h.errStr("failed to parse request", err))
				}
				reqByte = rByte
			}

			fields = append(fields, []logger.Field{
				{
					Key:   "response",
					Value: resBody,
				},
				{
					Key:   "request",
					Value: reqByte,
				},
			}...)
			logger.Info(lStr, fields...)
		} else {
			logger.Info(lStr, fields...)
		}
	}
}

func (h *HttpClient) errorLogging(method, url string, request io.Reader, start time.Time, err error, fields ...logger.Field) {

	l := time.Since(start).Round(time.Millisecond).String()
	lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s | %s", method,
		url, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), l)

	var reqByte []byte
	if request != nil {
		rByte, err := io.ReadAll(request)
		if err != nil {
			logger.Error(h.errStr("failed to parse request", err))
		}
		reqByte = rByte
	}

	fields = append(fields, []logger.Field{
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
			Value: reqByte,
		},
	}...)

	logger.Error(lStr, fields...)
}

// Request responsible for sending http request
// by using the Option set at the time of initialization.
func (h *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string,
	request io.Reader) (*http.Response, error) {
	var reqBuffer bytes.Buffer
	var reqReader io.Reader

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

	if request != nil {
		reader := io.TeeReader(request, &reqBuffer)
		reqReader = reader
	}

	rBody, err := http.NewRequestWithContext(ctx, method, url, reqReader)
	if err != nil {
		h.errorLogging(method, url, reqReader, start, err, lBody...)
		return nil, err
	}

	for key, val := range headers {
		rBody.Header[key] = []string{val}
	}

	response, err := h.client.Do(rBody)
	if err != nil {
		h.errorLogging(method, url, &reqBuffer, start, err, lBody...)
		return nil, err
	}

	h.successLogging(method, url, &reqBuffer, response.Status, response.StatusCode, response.Body, start, lBody...)

	return response, nil
}

// TODO: SOAP REQUEST
