package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	categoriesURL = "https://reference.craigslist.org/Categories"
	locationsURL  = "http://reference.craigslist.org/Areas"
)

// API houses possible queries to craigslist
type API interface {
	ListCategories()
	PrintCategories()
	ListLocations()
	PrintLocations()
}

// Client represents the main entrypoint to the API
type Client struct {
	initialized bool
	Categories  map[string]Category
	Locations   map[string]Location
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

// SubArea is a representation of a searchable area within a location
// {"Abbreviation":"sfc","Description":"city of san francisco","ShortDescription":"san francisco","SubAreaID":1}
type SubArea struct {
	Abbreviation     string
	Description      string
	ShortDescription string
	SubAreaID        int
}

type locationJSON struct {
	result []Location
}

type subareaJSON struct {
	result []SubArea
}

type categoriesJSON struct {
	result []Category
}

// Initialize will instantiate datastructures on the Client
func (c *Client) Initialize() {
	c.initialized = true
	c.Categories = make(map[string]Category)
	c.Locations = make(map[string]Location)
}

// ListCategories provides a list of all categories available on craigslist for query
// Reference: https://www.craigslist.org/about/reference
func (c *Client) ListCategories() (map[string]Category, error) {
	if !c.initialized {
		c.Initialize()
	}

	resp, err := http.Get(categoriesURL)
	noErr(err)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("category list failed: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	noErr(err)

	JSONtoSlice := categoriesJSON{}
	err = json.Unmarshal(b, &JSONtoSlice.result)
	noErr(err)

	for _, cat := range JSONtoSlice.result {
		c.Categories[cat.Abbreviation] = cat
	}

	return c.Categories, nil
}

// PrintCategories prints a list of all valid categories
func (c *Client) PrintCategories() {
	fmt.Println("-- CATEGORIES")
	for key, obj := range c.Categories {
		fmt.Printf("\n\raccess key: %s, details: %+v", key, obj)
	}
}

// ListLocations provides a list of all locations available on craiglist for query
// Reference: https://www.craigslist.org/about/reference
func (c *Client) ListLocations() (map[string]Location, error) {
	if !c.initialized {
		c.Initialize()
	}

	resp, err := http.Get(locationsURL)
	noErr(err)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("category list failed: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	noErr(err)

	JSONtoSlice := locationJSON{}
	err = json.Unmarshal(b, &JSONtoSlice.result)
	noErr(err)

	for _, loc := range JSONtoSlice.result {
		c.Locations[loc.Abbreviation] = loc
	}

	return c.Locations, nil
}

// PrintLocations prints a list of all valid locations
func (c *Client) PrintLocations() {
	fmt.Println("-- LOCATIONS")
	for key, obj := range c.Locations {
		fmt.Printf("\n\raccess key: %s, details: %+v", key, obj)
	}
}
