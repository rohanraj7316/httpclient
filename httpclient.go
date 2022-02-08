package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rohanraj7316/logger"
	"go.uber.org/zap/zapcore"
)

type HttpClient struct {
	client            *http.Client
	reqResLogging     bool
	reqResBodyLogging bool
}

func NewHTTPClient(o Options) (*HttpClient, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// setting up proxy
	if o.UseProxy {
		pProxyURL, err := url.Parse(o.ProxyURL)
		if err != nil {
			return nil, err
		}

		transport.Proxy = http.ProxyURL(pProxyURL)
	}

	err := logger.Configure(o.LoggerOptions)
	if err != nil {
		return nil, err
	}

	return &HttpClient{
		client: &http.Client{
			Timeout:   o.Timeout,
			Transport: transport,
		},
		reqResLogging:     o.LogReqResEnable,
		reqResBodyLogging: o.LogReqResBodyEnable,
	}, nil
}

func (h *HttpClient) successLogging(method, url, status string, statusCode int, headers, request, response interface{}) {
	if h.reqResLogging {
		lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s", method, url, statusCode, status)
		if h.reqResBodyLogging {
			reqResLogger := []zapcore.Field{
				{
					Key:       "request",
					Type:      zapcore.ReflectType,
					Interface: request,
				},
				{
					Key:       "response",
					Type:      zapcore.ReflectType,
					Interface: response,
				},
				{
					Key:       "headers",
					Type:      zapcore.ReflectType,
					Interface: headers,
				},
			}
			logger.Info(lStr, reqResLogger...)
		} else {
			logger.Info(lStr)
		}
	}
}

func (h *HttpClient) errorLogging(method, url string, request, headers interface{}, err error) {
	lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s", method, url, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	reqResLogger := []zapcore.Field{
		{
			Key:    "error",
			Type:   zapcore.StringType,
			String: err.Error(),
		},
	}
	if h.reqResBodyLogging {
		reqResLogger = append(reqResLogger, zapcore.Field{
			Key:       "request",
			Type:      zapcore.ReflectType,
			Interface: request,
		}, zapcore.Field{
			Key:       "headers",
			Type:      zapcore.ReflectType,
			Interface: headers,
		})
		logger.Error(lStr, reqResLogger...)
	} else {
		logger.Error(lStr)
	}
}

// Request responsible for sending http request
// by using the Option set at the time of initialization.
func (h *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string,
	request io.Reader) (*http.Response, error) {
	rBody, err := http.NewRequestWithContext(ctx, method, url, request)
	if err != nil {
		h.errorLogging(method, url, request, headers, err)
		return nil, err
	}

	for key, val := range headers {
		rBody.Header[key] = []string{val}
	}

	response, err := h.client.Do(rBody)
	if err != nil {
		h.errorLogging(method, url, request, headers, err)
		return nil, err
	}

	h.successLogging(method, url, response.Status, response.StatusCode, headers, request, response.Body)

	return response, nil
}

// TODO: SOAP REQUEST
