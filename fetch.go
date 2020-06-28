package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type fetcher interface {
	fetch(ctx context.Context, url string) (io.ReadCloser, error)
}

type httpService struct{}

func newHTTPService() fetcher {
	return &httpService{}
}

// simple function to isolate http requests from other services
func (f *httpService) fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("error fetching from url: %s", resp.Status)
	}

	return resp.Body, nil
}
