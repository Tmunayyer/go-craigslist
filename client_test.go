package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCategories []Category

var testLocations []Location

// A function to initialize a client without making the http
// requests for the categories and locations.
func newMockClient() (*client, error) {
	c := client{}

	c.initialized = true
	c.Categories = make(map[string]Category)
	c.Locations = make(map[string]Location)

	return &c, nil
}

func createMockResponse(t *testing.T, body interface{}) *httptest.ResponseRecorder {
	mockResponse := httptest.NewRecorder()

	mockBody, err := json.Marshal(body)
	assert.NoError(t, err)

	mockResponse.Code = 200
	_, err = mockResponse.Write(mockBody)
	assert.NoError(t, err)

	return mockResponse
}

func setupTests(t *testing.T) {
	testCategories = []Category{
		{
			Abbreviation: "mtk",
			CategoryID:   1,
			Description:  "money take kong",
			Type:         "T",
		},
		{
			Abbreviation: "tem",
			CategoryID:   2,
			Description:  "thomas elias munayyer",
			Type:         "H",
		},
	}

	testLocations = []Location{
		{
			Abbreviation:     "sf",
			AreaID:           1,
			Country:          "USA",
			Description:      "san fransisco",
			Hostname:         "hostname",
			Latitude:         1.23,
			Longitude:        2.13,
			Region:           "california",
			ShortDescription: "san fran",
			Timezone:         "california",
		},
		{
			Abbreviation:     "nyc",
			AreaID:           2,
			Country:          "USA",
			Description:      "the big apple",
			Hostname:         "newyork",
			Latitude:         1.23,
			Longitude:        2.13,
			Region:           "new york",
			ShortDescription: "New York city",
			Timezone:         "East Coast",
		},
	}
}

func TestSetCategories(t *testing.T) {
	setupTests(t)
	// skip the actual call and make sure the setCategories function
	// is working and errors appropriatly

	mockClient, _ := newMockClient()
	mockResponse := createMockResponse(t, testCategories)

	t.Run("setting categories on client", func(t *testing.T) {
		_, err := setCategories(mockClient, mockResponse.Result())
		assert.NoError(t, err)

		for _, cat := range testCategories {
			_, has := mockClient.Categories[cat.Abbreviation]
			assert.True(t, has)
		}
	})
}

func TestSetLocations(t *testing.T) {
	setupTests(t)
	// skip the actual call and make sure the setCategories function
	// is working and errors appropriatly

	mockClient, _ := newMockClient()
	mockResponse := createMockResponse(t, testLocations)

	t.Run("setting locations on client", func(t *testing.T) {
		_, err := setLocations(mockClient, mockResponse.Result())
		assert.NoError(t, err)

		for _, loc := range testLocations {
			_, has := mockClient.Locations[loc.Abbreviation]
			assert.True(t, has)
		}
	})
}

func TestQueryBuilder(t *testing.T) {
	setupTests(t)

	mockClient, err := newMockClient()
	assert.NoError(t, err)

	locationsBody := createMockResponse(t, testLocations).Result()
	categoriesBody := createMockResponse(t, testCategories).Result()

	_, err = setLocations(mockClient, locationsBody)
	assert.NoError(t, err)
	_, err = setCategories(mockClient, categoriesBody)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		filter   Filters
		expected string
	}{
		{
			name:     "q builder just loc and cat",
			args:     []string{"nyc", "mtk", "", ""},
			filter:   Filters{},
			expected: "https://newyork.craigslist.org/d/placeholder/search/mtk",
		},
		{
			name:     "q builder with a search term",
			args:     []string{"nyc", "mtk", "hello world", ""},
			filter:   Filters{},
			expected: "https://newyork.craigslist.org/d/placeholder/search/mtk?query=hello+world",
		},
		{
			name:     "q builder with a search term",
			args:     []string{"nyc", "mtk", "hello world %$#", ""},
			filter:   Filters{},
			expected: "https://newyork.craigslist.org/d/placeholder/search/mtk?query=hello+world+%25%24%23",
		},
		{
			name: "q builder with filters",
			args: []string{"nyc", "mtk", "hello world %$#", ""},
			filter: Filters{
				srchType:         true,
				hasPic:           true,
				postedToday:      true,
				bundleDuplicates: true,
				searchNearby:     true,
			},
			expected: "https://newyork.craigslist.org/d/placeholder/search/mtk?query=hello+world+%25%24%23&srchType=T&hasPic=1&postedToday=1&bundleDuplicates=1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := mockClient.BuildQuery(test.args[0], test.args[1], test.args[2], test.filter)
			assert.NoError(t, err)

			assert.Equal(t, test.expected, actual.URL)
		})
	}
}
