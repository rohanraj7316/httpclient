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

// Request responsible for sending http request
// by using the Option set at the time of initialization.
func (h *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string,
	request io.Reader) (*http.Response, error) {
	rBody, err := http.NewRequestWithContext(ctx, method, url, request)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		rBody.Header[key] = []string{val}
	}

	response, err := h.client.Do(rBody)
	if err != nil {
		return nil, err
	}

	if h.reqResLogging {
		lStr := fmt.Sprintf("HttpClient | %s | %s | %d | %s", method, url, response.StatusCode, response.Status)
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
					Interface: response.Body,
				},
			}
			logger.Info(lStr, reqResLogger...)
		} else {
			logger.Info(lStr)
		}
	}
	return response, nil
}

// TODO: SOAP REQUEST
