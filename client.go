package gocraigslist

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
)

// Client represents the interface with Craigslist.
type Client interface {
	FormatURL(term string, options Options) string
	GetListings(ctx context.Context, url string) ([]Listing, error)
}

type client struct {
	location string
}

// Options represents available parameters to construct a URL. Filters
// with tuple values are represented as [input value, mapped value].
type Options struct {
	location          string   // OPTIONAL: defaults to location provided on intialization, providing location here will overrides init value
	category          string   // OPTIONAL: defaults to constant defCategory, providing category overrides default variable
	postedBy          string   // OPTIONAL: [all, sss], [owner, sso], [dealer, ssq] attention: this only works for default search (sss), not specific categories
	srchType          bool     // OPTIONAL: true or false; dev note - uses "T" or "F" instead of 1 or 0
	hasPic            bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	postedToday       bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	bundleDuplicates  bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	cryptoCurrencyOK  bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	deliveryAvailable bool     // OPTIONAL: true or false; dev note - uses 1 for true, 0 for false
	minPrice          string   // OPTIONAL: example = 100
	maxPrice          string   // OPTIONAL: example = 500
	condition         []string // OPTIONAL: [new, 10], [like new, 20], [excellent, 30], [good, 40], [fair, 50], [salvage, 60]
	language          []string // OPTIONAL: [af, 1], [ca, 2], [da, 3], [de, 4], [en, 5], [es, 6], [fi, 7], [fr, 8], [it, 9], [nl, 10], [no, 11], [pt, 12], [sv, 13], [tl, 14], [tr, 15], [zh, 16], [ar, 17], [ja, 18], [ko, 19], [ru, 20], [vi, 21]
}

// NewClient will instantiate a client, set the location, and return a pointer.
func NewClient(location string) Client {
	c := client{location: location}
	return &c
}

// FormatURL should be used to programatically construct a URL using a term and Options.
func (c *client) FormatURL(term string, options Options) string {
	finalLocation := options.location
	if finalLocation == "" {
		finalLocation = c.location
	}

	finalCategory := options.category
	if finalCategory == "" {
		if options.postedBy == "owner" {
			finalCategory = defCategoryOwner
		} else if options.postedBy == "dealer" {
			finalCategory = defCategoryDealer
		} else {
			finalCategory = defCategory
		}
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

	conditionMap := map[string]string{"new": "10", "like new": "20", "excellent": "30", "good": "40", "fair": "50", "salvage": "60"}
	if len(options.condition) > 0 {
		for _, c := range options.condition {
			url += condition + conditionMap[c]
		}
	}

	languageMap := map[string]string{"af": "1", "ca": "2", "da": "3", "de": "4", "en": "5", "es": "6", "fi": "7", "fr": "8", "it": "9", "nl": "10", "no": "11", "pt": "12", "sv": "13", "tl": "14", "tr": "15", "zh": "16", "ar": "17", "ja": "18", "ko": "19", "ru": "20", "vi": "21"}
	if len(options.language) > 0 {
		for _, l := range options.language {
			url += language + languageMap[l]
		}
	}

	url += args

	return url
}

// GetListings simply takes a URL and returns the first page of listings.
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

func formatTerm(term string) string {
	pieces := strings.Split(term, " ")

	escapedPieces := []string{}
	for _, piece := range pieces {
		escapedPieces = append(escapedPieces, url.QueryEscape(piece))
	}

	return strings.Join(escapedPieces, "+")
}
