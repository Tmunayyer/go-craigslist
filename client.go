// 1 interface
// 2 struct
// 3 constructor
// 4 methods

package main

import (
	"context"
	"fmt"
	"net/http"
)

// Client houses possible queries to craigslist
type Client interface {
	// Prerequisites
	Initialize(ctx context.Context, location string) error

	// Primary Methods
	FormatURL(term string, location string, category string, filters Filters) (string, error)
	Search(ctx context.Context, url string) ([]Listing, error)
}

// Client represents the main entrypoint to the API
type client struct {
	initialized bool
	location    string
}

const (
	// defaults
	base        = "craigslist.org"
	defCategory = "sss" // this represent a text search in craigslist and not a specific category
	defPath     = "/search/"
	protocol    = "https"

	// defined queries and options
	sort             = "&sort="
	srchType         = "&srchType="
	hasPic           = "&hasPic="
	postedToday      = "&postedToday="
	bundleDuplicates = "&bundleDuplicates="
)

// NewClient needs context for the Initialize function. Initialize will make two http requests
// in order to get defined categories and areas from craigslist. These properties populate maps
// to provide a library later and validation to prevent bad requests from Search.
func NewClient(ctx context.Context, location string) (Client, error) {
	c := client{}
	err := c.Initialize(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("error initializing clinet: %v", err)
	}

	return &c, nil
}

// Initialize will instantiate datastructures on the Client struct
func (c *client) Initialize(ctx context.Context, location string) error {
	c.location = location
	c.initialized = true

	return nil
}

// Filters represents available filters
type Filters struct {
	srchType         bool // options: (T) true or (F) false
	hasPic           bool // options: (1) true or (0) false
	postedToday      bool // options: (1) true or (0) false
	bundleDuplicates bool // options: (1) true or (0) false
}

// FormatURL should be used when programatically accessing multiple types of queries. If a known page is desired,
// its best to provide the url directly to the search function.
func (c *client) FormatURL(term string, location string, category string, filters Filters) (string, error) {
	finalLocation := location
	if finalLocation == "" {
		finalLocation = c.location
	}

	finalCategory := category
	if finalCategory == "" {
		finalCategory = defCategory
	}

	var url string
	url = protocol + "://" + finalLocation + "." + base + defPath + finalCategory + "?query=" + term

	var args string
	if filters.srchType {
		args += srchType + "T"
	}

	if filters.postedToday {
		args += postedToday + "1"
	}

	if filters.hasPic {
		args += hasPic + "1"
	}

	if filters.bundleDuplicates {
		args += bundleDuplicates + "1"
	}

	url += args

	return url, nil
}

// Search takes in a url and parses the HTML response from Craigslist outputting a slice of listings
func (c *client) Search(ctx context.Context, url string) ([]Listing, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("location fetching from query url: %s", resp.Status)
	}

	listings, err := parseSearchResults(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing search results: %v", err)
	}

	return listings, nil
}
