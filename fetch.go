package main

import (
	"context"
	"fmt"
	"net/http"
)

type fetcher interface {
	fetch(ctx context.Context, url string) (*http.Response, error)
}

type httpService struct{}

func newHTTPService() fetcher {
	return &httpService{}
}

// simple function to isolate http requests from other services
func (f *httpService) fetch(ctx context.Context, url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("error fetching from url: %s", resp.Status)
	}

	return resp, nil
}
