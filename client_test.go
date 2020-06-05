import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCategories(t *testing.T) {
	// skip the actual call and make sure the setCategories function
	// is working and errors appropriatly

	mockClient := client{}
	mockClient.Initialize()

	mockResponse := httptest.NewRecorder()
	testCategories := []Category{
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

	mockBody, err := json.Marshal(testCategories)
	assert.NoError(t, err)

	mockResponse.Code = 200
	_, err = mockResponse.Write(mockBody)
	assert.NoError(t, err)

	t.Run("setting categories on client", func(t *testing.T) {
		setCategories(&mockClient, mockResponse.Result())

		for _, cat := range testCategories {
			_, has := mockClient.Categories[cat.Abbreviation]
			assert.True(t, has)
		}
	})
}

func TestSetLocations(t *testing.T) {
	// skip the actual call and make sure the setCategories function
	// is working and errors appropriatly

	mockClient := Client{}
	mockClient.Initialize()

	mockResponse := httptest.NewRecorder()
	testLocations := []Location{
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
			Hostname:         "hostname",
			Latitude:         1.23,
			Longitude:        2.13,
			Region:           "new york",
			ShortDescription: "New York city",
			Timezone:         "East Coast",
		},
	}

	mockBody, err := json.Marshal(testLocations)
	assert.NoError(t, err)

	mockResponse.Code = 200
	_, err = mockResponse.Write(mockBody)
	assert.NoError(t, err)

	t.Run("setting locations on client", func(t *testing.T) {
		setLocations(&mockClient, mockResponse.Result())

		for _, loc := range testLocations {
			_, has := mockClient.Locations[loc.Abbreviation]
			assert.True(t, has)
		}
	})
}
