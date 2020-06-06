// 1 interface
// 2 struct
// 3 constructor
// 4 methods

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client houses possible queries to craigslist
type Client interface {
	// Primary Methods
	FetchCategories(ctx context.Context) (map[string]Category, error)
	FetchLocations(ctx context.Context) (map[string]Location, error)
	BuildQuery(loc string, cat string, term string, filters Filters) (Query, error)

	// Convenience
	PrintCategories() // Prints categories on the client
	PrintLocations()  // Prints locations on the client
}

// Client represents the main entrypoint to the API
type client struct {
	initialized bool
	Categories  map[string]Category
	Locations   map[string]Location
}

const (
	categoriesURL = "https://reference.craigslist.org/Categories"
	locationsURL  = "http://reference.craigslist.org/Areas"
)

// NewClient needs context for the Initialize function. Initialize will make two http requests
// in order to get defined categories and areas from craigslist. These properties populate maps
// to provide a library later and validation to prevent bad requests from Search.
func NewClient(ctx context.Context) (Client, error) {
	c := client{}
	err := c.Initialize(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing clinet: %v", err)
	}

	return &c, nil
}

// Initialize will instantiate datastructures on the Client struct
func (c *client) Initialize(ctx context.Context) error {
	// These functions are required to run before any other functions can be run. The purpose
	// of these functions is to get an up to date mapping of allowed locations and categories
	// from craigslist. These value will populate maps on the client that will later be used
	// for query validation in BuildQuery to prevent errored responses or invalid parameters
	// for the Search function

	c.initialized = true
	c.Categories = make(map[string]Category)
	c.Locations = make(map[string]Location)

	_, err := c.FetchLocations(ctx)
	if err != nil {
		return fmt.Errorf("error fetching locations: %v", err)
	}
	_, err = c.FetchCategories(ctx)
	if err != nil {
		return fmt.Errorf("error fetching categories: %v", err)
	}

	return nil
}

// Category represents a valid Craigslist category that can be queried
type Category struct {
	Abbreviation string
	CategoryID   int
	Description  string
	Type         string
}

// Location represents a single valid craigslist location
type Location struct {
	Abbreviation     string
	AreaID           int
	Country          string
	Description      string
	Hostname         string
	Latitude         float32
	Longitude        float32
	Region           string
	ShortDescription string
	Timezone         string
	SubAreas         []SubArea
}

// SubArea is a representation of a searchable area within a location, these are
// filters that can be provided in a given search using the SubAreaID
type SubArea struct {
	Abbreviation     string
	Description      string
	ShortDescription string
	SubAreaID        int
}

// FetchLocations provides a list of all locations available on craiglist for query
// Reference: https://www.craigslist.org/about/reference
func (c *client) FetchLocations(ctx context.Context) (map[string]Location, error) {
	resp, err := http.Get(locationsURL)
	if err != nil {
		return nil, fmt.Errorf("error getting locations: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("location list failed: %s", resp.Status)
	}

	return setLocations(c, resp)
}

func setLocations(c *client, resp *http.Response) (map[string]Location, error) {
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	JSONtoSlice := []Location{}
	err = json.Unmarshal(b, &JSONtoSlice)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %v", err)
	}

	for _, loc := range JSONtoSlice {
		c.Locations[loc.Abbreviation] = loc
	}

	return c.Locations, nil
}

// FetchCategories provides a list of all categories available on craigslist for query
// Reference: https://www.craigslist.org/about/reference
func (c *client) FetchCategories(ctx context.Context) (map[string]Category, error) {
	resp, err := http.Get(categoriesURL)
	if err != nil {
		return nil, fmt.Errorf("error getting categories: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("category list failed: %s", resp.Status)
	}

	return setCategories(c, resp)
}

func setCategories(c *client, resp *http.Response) (map[string]Category, error) {
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	JSONtoSlice := []Category{}
	err = json.Unmarshal(b, &JSONtoSlice)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %v", err)
	}

	for _, cat := range JSONtoSlice {
		c.Categories[cat.Abbreviation] = cat
	}

	return c.Categories, nil
}

// PrintCategories prints a list of all valid categories on the Client struct
func (c *client) PrintCategories() {
	fmt.Println("-- CATEGORIES")
	for key, obj := range c.Categories {
		fmt.Printf("\n\raccess key: %s, details: %+v", key, obj)
	}
}

// PrintLocations prints a list of all valid locations on the Client struct
func (c *client) PrintLocations() {
	fmt.Println("-- LOCATIONS")
	for key, obj := range c.Locations {
		fmt.Printf("\n\raccess key: %s, details: %+v", key, obj)
	}
}

// Filters represents available filters
type Filters struct {
	srchType         bool // options: (T) true or (F) false
	hasPic           int  // options: (1) true or (0) false
	postedToday      int  // options: (1) true or (0) false
	bundleDuplicates int  // options: (1) true or (0) false
}

// Query is the struct representation of a query passed to the Search function
type Query struct {
	Location string
	Category string
	Term     string
	Filters  string
}

// BuildQuery will take in its arguments and format a query pattern
// TODO: Implement filters
func (c *client) BuildQuery(loc string, cat string, term string, filters Filters) (Query, error) {
	// validate location
	_, has := c.Locations[loc]
	if !has {
		return Query{}, fmt.Errorf("invalid location provided: %s", loc)
	}

	// validate category
	_, has = c.Categories[cat]
	if !has {
		return Query{}, fmt.Errorf("invalid category provided: %s", cat)
	}

	// build the term
	// TODO: ensure this is replacing " " with "+"
	escapedTerm := url.PathEscape(term)

	return Query{
		Location: loc,
		Category: cat,
		Term:     escapedTerm,
	}, nil
}
