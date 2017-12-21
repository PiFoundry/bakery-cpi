package bakeryclient

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	url        string
}

func New(url string) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
