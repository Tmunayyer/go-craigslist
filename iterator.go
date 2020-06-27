package main

import (
	"context"
	"fmt"
	"strconv"
)

// Iterator is an iterator used to retrieve multiple pages of listings
type Iterator interface {
	Next(context.Context) (*Result, error)
}

// Result is the structural representation of Result
type Result struct {
	Client      API
	Done        bool
	Listings    []Listing
	TotalCount  int
	CurrentPage int
	SearchURL   string // the original search url without pagination
	Err         error
}

func newResult(c API, url string, totalCount int, listings []Listing) *Result {
	r := Result{
		Client:      c,
		Done:        false,
		Listings:    listings,
		TotalCount:  totalCount,
		CurrentPage: 0,
		SearchURL:   url,
	}

	if len(listings) >= totalCount {
		r.Done = true
	}

	return &r
}

// Next page of listings. This will call the librarys fn GetListings.
// Note: Because postings could be happening as this is fetching results, there is
// the possibility of some duplicates coming through.
// 		Example in seconds:
//			time 0: 1st page is fetched
//			time 1: 1st page of listings is being processed, new listing is posted
// 			time 2: 2nd page is fetched
//		the second page would contain the last listing of the previous page.
func (r *Result) Next(ctx context.Context) (*Result, error) {
	r.CurrentPage++
	nextPageStart := r.CurrentPage * 120
	nextPageURL := r.SearchURL + page + strconv.Itoa(nextPageStart)

	respBody, err := fetch(ctx, nextPageURL)
	if err != nil {
		respBody.Close()
		r.Done = true
		r.Listings = []Listing{}
		return r, fmt.Errorf("error fetching from url: %v", err)
	}

	listings, _, err := parseSearchResults(respBody)
	if err != nil {
		r.Done = true
		return r, err
	}

	r.Listings = listings
	if (nextPageStart + len(r.Listings)) >= r.TotalCount {
		r.Done = true
	}

	return r, nil
}
