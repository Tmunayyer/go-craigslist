// 1 interface
// 2 struct
// 3 constructor
// 4 methods

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client houses possible queries to craigslist
type Client interface {
	// Prerequisites
	Initialize(ctx context.Context, location string) error

	// Primary Methods
	FormatURL(term string, options Options) (string, error)
	GetListings(ctx context.Context, url string) ([]Listing, error)
}

// Client represents the main entrypoint to the API
type client struct {
	initialized bool
	location    string
}

const (
	// defaults
	protocol    = "https"
	base        = "craigslist.org"
	defCategory = "sss" // this represent a text search in craigslist and not a specific category
	defPath     = "/search/"
	defSort     = "&sort=rel"

	// other useful constants
	defCategoryOwner  = "sso"
	defCategoryDealer = "ssq"

	// defined queries and options

	srchType          = "&srchType="
	hasPic            = "&hasPic="
	postedToday       = "&postedToday="
	bundleDuplicates  = "&bundleDuplicates="
	cryptoCurrencyOK  = "&crypto_currency_ok="
	deliveryAvailable = "&delivery_available="
	minPrice          = "&min_price="
	maxPrice          = "&max_price="
	language          = "&language="
	condition         = "&condition="
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

// Options represents available filters
type Options struct {
	// params
	location string // OPTIONAL: defaults to location provided on intialization, providing location here will overrides init value
	category string // OPTIONAL: defaults to constant defCategory, providing category overrides default variable

	// filters
	postedBy          int      // OPTIONAL: [all, 1], [owner, 2], [dealder, 3] note: this only works for general search, not specific categories
	srchType          bool     // OPTIONAL: dev note - uses "T" or "F" instead of 1 or 0
	hasPic            bool     // OPTIONAL: dev note - uses 1 for true, 0 for false
	postedToday       bool     // OPTIONAL: dev note - uses 1 for true, 0 for false
	bundleDuplicates  bool     // OPTIONAL: dev note - uses 1 for true, 0 for false
	cryptoCurrencyOK  bool     // OPTIONAL: dev note - uses 1 for true, 0 for false
	deliveryAvailable bool     // OPTIONAL: dev note - uses 1 for true, 0 for false
	minPrice          string   // OPTIONAL: example = 100
	maxPrice          string   // OPTIONAL: example = 500
	condition         []string // OPTIONAL: [new, 10], [like new, 20], [excellent, 30], [good, 40], [fair, 50], [salvage, 60]
	language          []string // OPTIONAL: [af, 1], [ca, 2], [da, 3], [de, 4], [en, 5], [es, 6], [fi, 7], [fr, 8], [it, 9], [nl, 10], [no, 11], [pt, 12], [sv, 13], [tl, 14], [tr, 15], [zh, 16], [ar, 17], [ja, 18], [ko, 19], [ru, 20], [vi, 21]
}

// FormatURL is used for programaticaly constructing craigslist search urls.
func (c *client) FormatURL(term string, options Options) (string, error) {
	// fmt.Printf("the options: %+v", options)
	finalLocation := options.location
	if finalLocation == "" {
		finalLocation = c.location
	}

	finalCategory := options.category
	if finalCategory == "" {
		finalCategory = defCategory
	}

	formattedTerm := formatTerm(term)

	var url string
	url = protocol + "://" + finalLocation + "." + base + defPath + finalCategory + "?query=" + formattedTerm + defSort

	var args string
	if options.srchType {
		args += srchType + "T"
	}

	if options.postedToday {
		args += postedToday + "1"
	}

	if options.hasPic {
		args += hasPic + "1"
	}

	if options.bundleDuplicates {
		args += bundleDuplicates + "1"
	}

	if options.cryptoCurrencyOK {
		args += cryptoCurrencyOK + "1"
	}

	if options.deliveryAvailable {
		args += deliveryAvailable + "1"
	}

	if options.minPrice != "" {
		args += minPrice + options.minPrice
	}

	if options.maxPrice != "" {
		args += maxPrice + options.maxPrice
	}

	url += args

	return url, nil
}

func formatTerm(term string) string {
	pieces := strings.Split(term, " ")

	escapedPieces := []string{}
	for _, piece := range pieces {
		escapedPieces = append(escapedPieces, url.QueryEscape(piece))
	}

	return strings.Join(escapedPieces, "+")
}

// GetListings simply takes a URL and returns the first page of listings on that page.
func (c *client) GetListings(ctx context.Context, url string) ([]Listing, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("error fetching from url: %s", resp.Status)
	}

	listings, err := parseSearchResults(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing search results: %v", err)
	}

	return listings, nil
}
