package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func setCtx(pCtx context.Context) context.Context {
	nId, _ := uuid.NewUUID()
	return context.WithValue(pCtx, "requestId", nId)
}

func TestHttpMethods(t *testing.T) {
	ctx := context.Background()

	// add your test case here
	tests := []struct {
		ctx     context.Context
		url     string
		method  string
		headers map[string]string
		request io.Reader
	}{
		{
			ctx:    setCtx(ctx),
			url:    "https://httpbin.org/anything",
			method: http.MethodGet,
			headers: map[string]string{
				"T-ID": "123123",
			},
			request: nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing %s Method", tt.method), func(t *testing.T) {
			// intializing the request
			c, _ := New()

			_, err := c.Request(tt.ctx, tt.method, tt.url, tt.headers, tt.request)
			if err != nil {
				t.Error(err)
			}

			t.Log(fmt.Sprintf("successfully tested %s", tt.method))
		})
	}
}

func TestFlags(t *testing.T) {
	ctx := context.Background()

	tests := []Config{
		{
			LogReqResEnable:     true,
			LogReqResBodyEnable: true,
		},
		{
			LogReqResEnable:     true,
			LogReqResBodyEnable: false,
		},
		{
			LogReqResEnable:     false,
			LogReqResBodyEnable: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing Flag LogReqResEnable: %t & LogReqResBodyEnable: %t", tt.LogReqResEnable, tt.LogReqResBodyEnable), func(t *testing.T) {
			cfg := Config{
				LogReqResEnable:     tt.LogReqResEnable,
				LogReqResBodyEnable: tt.LogReqResBodyEnable,
			}
			c, err := New(cfg)
			if err != nil {
				t.Error(err)
			}

			url := "https://httpbin.org/anything"
			headers := map[string]string{
				"T-ID": "123123",
			}
			_, err = c.Request(setCtx(ctx), http.MethodGet, url, headers, nil)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRequestLogging(t *testing.T) {
	ctx := context.Background()
	url := "https://jsonplaceholder.typicode.com/posts"

	c, err := New()
	if err != nil {
		t.Error(err)
	}

	_, err := c.Request(setCtx(ctx), http.MethodGet, url, nil, nil)
	if err != nil {
		t.Error(err)
	}
}
