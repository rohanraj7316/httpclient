package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rohanraj7316/httpclient"
)

type Handler struct {
	client *httpclient.HttpClient
}

func NewHandler() *Handler {
	client, err := httpclient.New()
	if err != nil {
		log.Println(err)
	}

	return &Handler{
		client: client,
	}
}

func (h *Handler) Get(ctx context.Context) {
	url := "https://httpbin.org/anything"
	header := map[string]string{
		"content-type": "application/json",
	}
	_, err := h.client.Request(ctx, http.MethodGet, url, header, nil)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) Post(ctx context.Context) {
	url := "https://httpbin.org/anything"
	header := map[string]string{
		"Content-Type": "application/json",
	}

	body := map[string]string{
		"name": "Rohan Raj",
	}
	bBody, err := json.Marshal(&body)
	if err != nil {
		log.Fatal(err)
	}

	_, err = h.client.Request(ctx, http.MethodPost, url, header, bytes.NewBuffer(bBody))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	nH := NewHandler()
	c := context.Background()
	ctx := context.WithValue(c, "requestId", "123123123123123")

	// GET
	nH.Get(ctx)

	// Post
	nH.Post(ctx)
}
