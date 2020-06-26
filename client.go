package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
	page              = "&s="
)

// API represents the interface with Craigslist.
type API interface {
	FormatURL(term string, options Options) string
	GetListings(ctx context.Context, url string) ([]Listing, error)
	GetNewListings(ctx context.Context, url string, date time.Time) ([]Listing, error)
}

// Client is return from New Client with a Location. This Location is used as
// the default value in FormatURL unless one is provided in Options.
type Client struct {
	Location string
}

// Options represents available parameters to construct a URL. Filters
// with tuple values are represented as [input value, mapped value].
type Options struct {
	Location          string   // OPTIONAL: defaults to location provided on intialization, providing location here will overrides init value
	Category          string   // OPTIONAL: defaults to constant defCategory, providing category overrides default variable
	PostedBy          string   // OPTIONAL: [all, sss], [owner, sso], [dealer, ssq] attention: this only works for default search (sss), not specific categories
	SrchType          bool     // OPTIONAL: true or false; dev note - uses "T" or "F" instead of 1 or 0
	HasPic            bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	PostedToday       bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	BundleDuplicates  bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	CryptoCurrencyOK  bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	DeliveryAvailable bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	MinPrice          string   // OPTIONAL: example = 100
	MaxPrice          string   // OPTIONAL: example = 500
	Condition         []string // OPTIONAL: [new, 10], [like new, 20], [excellent, 30], [good, 40], [fair, 50], [salvage, 60]
	Language          []string // OPTIONAL: [af, 1], [ca, 2], [da, 3], [de, 4], [en, 5], [es, 6], [fi, 7], [fr, 8], [it, 9], [nl, 10], [no, 11], [pt, 12], [sv, 13], [tl, 14], [tr, 15], [zh, 16], [ar, 17], [ja, 18], [ko, 19], [ru, 20], [vi, 21]
}

// NewClient will instantiate a client, set the location, and return a pointer.
func NewClient(location string) API {
	c := Client{Location: location}
	return &c
}

// FormatURL should be used to programatically construct a URL using a term and Options.
func (c *Client) FormatURL(term string, options Options) string {
	finalLocation := options.Location
	if finalLocation == "" {
		finalLocation = c.Location
	}

	finalCategory := options.Category
	if finalCategory == "" {
		if options.PostedBy == "owner" {
			finalCategory = defCategoryOwner
		} else if options.PostedBy == "dealer" {
			finalCategory = defCategoryDealer
		} else {
			finalCategory = defCategory
		}
	}

	formattedTerm := formatTerm(term)

	var url string
	url = protocol + "://" + finalLocation + "." + base + defPath + finalCategory + "?query=" + formattedTerm + defSort

	var args string
	if options.SrchType {
		args += srchType + "T"
	}

	if options.PostedToday {
		args += postedToday + "1"
	}

	if options.HasPic {
		args += hasPic + "1"
	}

	if options.BundleDuplicates {
		args += bundleDuplicates + "1"
	}

	if options.CryptoCurrencyOK {
		args += cryptoCurrencyOK + "1"
	}

	if options.DeliveryAvailable {
		args += deliveryAvailable + "1"
	}

	if options.MinPrice != "" {
		args += minPrice + options.MinPrice
	}

	if options.MaxPrice != "" {
		args += maxPrice + options.MaxPrice
	}

	conditionMap := map[string]string{"new": "10", "like new": "20", "excellent": "30", "good": "40", "fair": "50", "salvage": "60"}
	if len(options.Condition) > 0 {
		for _, c := range options.Condition {
			url += condition + conditionMap[c]
		}
	}

	languageMap := map[string]string{"af": "1", "ca": "2", "da": "3", "de": "4", "en": "5", "es": "6", "fi": "7", "fr": "8", "it": "9", "nl": "10", "no": "11", "pt": "12", "sv": "13", "tl": "14", "tr": "15", "zh": "16", "ar": "17", "ja": "18", "ko": "19", "ru": "20", "vi": "21"}
	if len(options.Language) > 0 {
		for _, l := range options.Language {
			url += language + languageMap[l]
		}
	}

	url += args

	return url
}

// GetListings simply takes a URL and returns the first page of listings.
func (c *Client) GetListings(ctx context.Context, url string) ([]Listing, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("error send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("error fetching from url: %s", resp.Status)
	}

	listings, count, err := parseSearchResults(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing search results: %v", err)
	}

	return listings, count, nil
}

// GetNewListings performs the same tasks as GetListings but only
// returns listings greater than the passed in date
func (c *Client) GetNewListings(ctx context.Context, url string, date time.Time) ([]Listing, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("error fetching from url: %s", resp.Status)
	}

	listings, err := parseSearchResultsAfter(resp.Body, date)
	if err != nil {
		return nil, fmt.Errorf("error parsing search results: %v", err)
	}

	return listings, nil
}

// GetMultipageListings returns an iterator to retrieve multiple pages of listings
func (c *Client) GetMultipageListings(ctx context.Context, url string) (Result, error) {
	// get the first page
	listings, count, err := c.GetListings(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error send request: %v", err)
	}

	r := newResult(c, url, count, listings)

	return r, nil
}

func formatTerm(term string) string {
	pieces := strings.Split(term, " ")

	escapedPieces := []string{}
	for _, piece := range pieces {
		escapedPieces = append(escapedPieces, url.QueryEscape(piece))
	}

	return strings.Join(escapedPieces, "+")
}
