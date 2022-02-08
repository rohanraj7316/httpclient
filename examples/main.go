package main

import (
	"context"
	"libs/httpclient"
	"log"
	"net/http"
	"time"

	"github.com/rohanraj7316/logger"
)

func NewLogConfig(o *logger.Options) (*logger.Options, error) {
	o.JSONEncoding = true
	o.IncludeCallerSourceLocation = true
	o.LogGrpc = true
	return o, nil
}

func main() {
	lOptions, _ := NewLogConfig(logger.NewOptions())
	hOptions := httpclient.Options{
		Timeout:             20 * time.Second,
		LoggerOptions:       lOptions,
		LogReqResEnable:     true,
		LogReqResBodyEnable: true,
	}

	client, err := httpclient.NewHTTPClient(hOptions)
	if err != nil {
		log.Println(err)
	}

	// GET
	ctx := context.Background()
	url := "https://httpbin.org/anything"
	header := map[string]string{
		"content-type": "application/json",
	}

	_, err = client.Request(ctx, http.MethodGet, url, header, nil)
	if err != nil {
		log.Println(err)
	}
}
