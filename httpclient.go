package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
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

// Request responsible for sending http request
// by using the Option set at the time of initialization.
func (h *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string,
	request io.Reader) (res *http.Response, err error) {

	var reqByte, resByte []byte
	start := time.Now()

	if request != nil {
		reqByte, err = ioutil.ReadAll(request)
		if err != nil {
			go h.errorLogging(ctx, method, url, reqByte, start, err)
			return nil, err
		}
	}

	rBody, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqByte))
	if err != nil {
		go h.errorLogging(ctx, method, url, reqByte, start, err)
		return nil, err
	}

	for key, val := range headers {
		rBody.Header[key] = []string{val}
	}

	res, err = h.client.Do(rBody)
	if err != nil {
		go h.errorLogging(ctx, method, url, reqByte, start, err)
		return nil, err
	}

	resByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		go h.errorLogging(ctx, method, url, reqByte, start, err)
		return nil, err
	}
	res.Body.Close()

	res.Body = ioutil.NopCloser(bytes.NewReader(resByte))

	go h.successLogging(ctx, method, url, res.Status, res.StatusCode, reqByte, resByte, start)

	return res, nil
}

// TODO: SOAP REQUEST
