package main

import (
	"context"
	"strconv"
)

// Result is an iterator used to retrieve multiple pages of listings
type Result interface {
	Next(context.Context) Result
}

// Iterator is the structural representation of Result
type Iterator struct {
	Client      *API
	Done        bool
	Listings    []Listing
	TotalCount  int
	CurrentPage int
	SearchURL   string // the original search url without pagination
	Err         error
}

func newResult(c *API, url string, totalCount int, listings []Listing) Result {
	i := Iterator{
		Client:      c,
		Done:        false,
		Listings:    listings,
		TotalCount:  totalCount,
		CurrentPage: 0,
		SearchURL:   url,
		Err:         nil,
	}

	if len(listings) <= totalCount {
		i.Done = true
	}

	return &i
}

// Next page of listings. This will call the librarys fn GetListings.
func (i *Iterator) Next(ctx context.Context) Result {
	i.CurrentPage++
	nextPageStart := i.CurrentPage * 120
	nextPageURL := i.SearchURL + page + strconv.Itoa(nextPageStart)

	listings, _, err := i.Client.GetListings(ctx, nextPageURL)
	if err != nil {
		i.Done = true
		i.Err = err
		return i
	}

	i.Listings = listings

	if (nextPageStart + 120) >= i.TotalCount {
		i.Done = true
	}

	return i
}
