package gocraigslist

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestFindBy(t *testing.T) {
	data, err := ioutil.ReadFile("./test.html")
	assert.NoError(t, err)

	r := bytes.NewReader(data)

	doc, err := html.Parse(r)
	assert.NoError(t, err)

	for _, test := range []struct {
		name     string
		inputs   []string
		expected bool
	}{
		{
			name:     "should find node",
			inputs:   []string{"id", "sortable-results"},
			expected: true,
		},
		{
			name:     "should not find node",
			inputs:   []string{"id", "thomas"},
			expected: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, has := findBy(doc, test.inputs[0], test.inputs[1])
			assert.Equal(t, test.expected, has)
		})
	}
}

func TestFindAttr(t *testing.T) {
	data, err := ioutil.ReadFile("./test.html")
	assert.NoError(t, err)

	r := bytes.NewReader(data)

	doc, err := html.Parse(r)
	assert.NoError(t, err)

	for _, test := range []struct {
		name     string
		inputs   string
		expected string
	}{
		{
			name:     "should find data",
			inputs:   "datetime",
			expected: "2020-06-08 14:45",
		},
		{
			name:     "should not find anything",
			inputs:   "yellow",
			expected: "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			node, _ := findBy(doc, "class", "result-date")
			_, data := findAttr(node.Attr, test.inputs)
			assert.Equal(t, test.expected, data)
		})
	}
}

func TestExtractListings(t *testing.T) {

	t.Run("no cutoff time provided", func(t *testing.T) {
		data, err := ioutil.ReadFile("./test.html")
		assert.NoError(t, err)

		r := bytes.NewReader(data)

		doc, err := html.Parse(r)
		assert.NoError(t, err)
		// find the entrypoint to  the results section of the page
		resultSection, _ := findBy(doc, "id", "sortable-results")
		// find the resultList, everything in here will go into the listing slice
		resultList, _ := findBy(resultSection, "class", "rows")
		listings := extractListings(resultList, nilTime)

		assert.Equal(t, 120, len(listings))
	})

	t.Run("cutoff date provided", func(t *testing.T) {
		data, err := ioutil.ReadFile("./test.html")
		assert.NoError(t, err)

		r := bytes.NewReader(data)

		doc, err := html.Parse(r)
		assert.NoError(t, err)
		// find the entrypoint to  the results section of the page
		resultSection, _ := findBy(doc, "id", "sortable-results")
		// find the resultList, everything in here will go into the listing slice
		resultList, _ := findBy(resultSection, "class", "rows")

		layout := "2006-01-02 15:04"
		cutoff, err := time.Parse(layout, "2020-06-08 14:03")
		assert.NoError(t, err)

		listings := extractListings(resultList, cutoff)

		assert.Equal(t, 19, len(listings))
	})

}
