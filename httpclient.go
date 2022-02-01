package httpclient

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	client *http.Client
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

	return &HttpClient{
		client: &http.Client{
			Timeout:   o.Timeout,
			Transport: transport,
		},
	}, nil
}

// Request responsible for sending http request
// by using the Option set at the time of initialization.
func (h *HttpClient) Request(ctx context.Context, method, url string, headers map[string]string,
	body io.Reader) (*http.Response, error) {
	// TODO: logging the request
	rBody, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		rBody.Header[key] = []string{val}
	}

	r, err := h.client.Do(rBody)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// TODO: SOAP REQUEST
