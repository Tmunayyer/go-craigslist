package main

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatURL(t *testing.T) {
	// note: expected variables should be pulled directly from craigslist
	// for all tests except randomized

	// this is a map of interfaces to aggregate all the
	// different types of possibilities resulting from options
	// argMap := struct {
	// 	truthy map[string]string
	// 	falsey map[string]string
	// }{
	// 	truthy: map[string]string{
	// 		// BOOLEANS
	// 		"srchType":          "T",
	// 		"hasPic":            "1",
	// 		"postedToday":       "1",
	// 		"bundleDuplicated":  "1",
	// 		"cryptoCurrencyOK":  "1",
	// 		"deliveryAvailable": "1",

	// 		// PRICES
	// 		// minPrice and maxPrice are variable, just check against actual options struct

	// 		// CONDITIONS
	// 		// language and conditions should be present if passed
	// 		// and should map to correct integer in string value
	// 		"new":       "10",
	// 		"like new":  "20",
	// 		"excellent": "30",
	// 		"good":      "40",
	// 		"fair":      "50",
	// 		"salvage":   "60",

	// 		// LANGUAGES
	// 		"af": "1",
	// 		"ca": "2",
	// 		"da": "3",
	// 		"de": "4",
	// 		"en": "5",
	// 		"es": "6",
	// 		"fi": "7",
	// 		"fr": "8",
	// 		"it": "9",
	// 		"nl": "10",
	// 		"no": "11",
	// 		"pt": "12",
	// 		"sv": "13",
	// 		"tl": "14",
	// 		"tr": "15",
	// 		"zh": "16",
	// 		"ar": "17",
	// 		"ja": "18",
	// 		"ko": "19",
	// 		"ru": "20",
	// 		"vi": "21",
	// 	},
	// 	falsey: map[string]string{
	// 		"srchType":          "F",
	// 		"hasPic":            "0",
	// 		"postedToday":       "0",
	// 		"bundleDuplicated":  "0",
	// 		"cryptoCurrencyOK":  "0",
	// 		"deliveryAvailable": "0",

	// 		// minPrice, maxPrice, language, and conditions should be missing if not passed
	// 	},
	// }

	client, err := NewClient(context.Background(), "newyork")
	assert.NoError(t, err)

	t.Run("no options provided, basic term", func(t *testing.T) {
		expected := "https://newyork.craigslist.org/search/sss?query=xbox&sort=rel"
		url, err := client.FormatURL("xbox", Options{})
		assert.NoError(t, err)
		assert.Equal(t, expected, url)
	})

	t.Run("properly escaping term", func(t *testing.T) {
		expected := "https://newyork.craigslist.org/search/sss?query=xbox+123+.+%23%24%25&sort=rel"
		url, err := client.FormatURL("xbox 123 . #$%", Options{})
		assert.NoError(t, err)
		assert.Equal(t, expected, url)
	})

	t.Run("should overwright location if provided by options", func(t *testing.T) {
		o := Options{location: "testing"}
		expected := "https://testing.craigslist.org/search/sss?query=xbox&sort=rel"
		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)
		assert.Equal(t, expected, url)
	})

	t.Run("should overwright category if provided by options", func(t *testing.T) {
		o := Options{category: "aaa"}
		expected := "https://newyork.craigslist.org/search/aaa?query=xbox&sort=rel"

		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)
		assert.Equal(t, expected, url)
	})

	t.Run("should account for bool options", func(t *testing.T) {
		o := Options{
			srchType:          true,
			hasPic:            true,
			postedToday:       true,
			bundleDuplicates:  true,
			cryptoCurrencyOK:  true,
			deliveryAvailable: true,
		}

		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)

		// lookup each arg by name in a map and ensure everything is there
		// with appropriate value

		beginQ := strings.Index(url, "?")
		urlArgs := strings.Split(url[beginQ+1:], "&")
		argMap := map[string]string{}
		for _, arg := range urlArgs {
			pieces := strings.Split(arg, "=")
			argMap[pieces[0]] = pieces[1]
		}

		numberOfOptions := 0
		if o.srchType {
			assert.Equal(t, "T", argMap["srchType"])
			numberOfOptions++
		}

		if o.hasPic {
			assert.Equal(t, "1", argMap["hasPic"])
			numberOfOptions++
		}

		if o.postedToday {
			assert.Equal(t, "1", argMap["postedToday"])
			numberOfOptions++
		}

		if o.postedToday {
			assert.Equal(t, "1", argMap["postedToday"])
			numberOfOptions++
		}

		if o.bundleDuplicates {
			assert.Equal(t, "1", argMap["bundleDuplicates"])
			numberOfOptions++
		}

		if o.cryptoCurrencyOK {
			assert.Equal(t, "1", argMap["crypto_currency_ok"])
			numberOfOptions++
		}

		if o.deliveryAvailable {
			assert.Equal(t, "1", argMap["delivery_available"])
			numberOfOptions++
		}

		// the query and rel arg are passed always, minus 2 from length of argMap
		assert.Equal(t, numberOfOptions, len(argMap)-1)
	})

	t.Run("accounts for min and max price", func(t *testing.T) {
		o := Options{
			minPrice: "100",
			maxPrice: "500",
		}

		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)

		beginQ := strings.Index(url, "?")
		urlArgs := strings.Split(url[beginQ+1:], "&")
		argMap := map[string]string{}
		for _, arg := range urlArgs {
			pieces := strings.Split(arg, "=")
			argMap[pieces[0]] = pieces[1]
		}

		assert.Equal(t, o.minPrice, argMap["min_price"])
		assert.Equal(t, o.maxPrice, argMap["max_price"])
	})

	t.Run("accounts for conditions", func(t *testing.T) {
		for _, test := range []struct {
			given    Options
			expected int
		}{
			{
				given:    Options{condition: []string{"new"}},
				expected: 10,
			},
			{
				given:    Options{condition: []string{"like new"}},
				expected: 20,
			},
			{
				given:    Options{condition: []string{"excellent"}},
				expected: 30,
			},
			{
				given:    Options{condition: []string{"good"}},
				expected: 40,
			},
			{
				given:    Options{condition: []string{"fair"}},
				expected: 50,
			},
			{
				given:    Options{condition: []string{"salvage"}},
				expected: 60,
			},
			{
				given:    Options{condition: []string{"new", "like new", "excellent", "good", "fair", "salvage"}},
				expected: 210, // 10 + 20 + 30 + 40 + 50 + 60
			},
		} {
			url, err := client.FormatURL("xbox", test.given)
			assert.NoError(t, err)

			analyzeURL(t, url, test.given, test.expected)
		}
	})

	t.Run("accounts for languages", func(t *testing.T) {
		languageMap := map[string]string{"af": "1", "ca": "2", "da": "3", "de": "4", "en": "5", "es": "6", "fi": "7", "fr": "8", "it": "9", "nl": "10", "no": "11", "pt": "12", "sv": "13", "tl": "14", "tr": "15", "zh": "16", "ar": "17", "ja": "18", "ko": "19", "ru": "20", "vi": "21"}

		allLanguages := []string{}
		for k := range languageMap {
			// push to a slice to use for testing all langs later
			allLanguages = append(allLanguages, k)

			o := Options{language: []string{k}}

			url, err := client.FormatURL("xbox", o)
			assert.NoError(t, err)

			total, err := strconv.Atoi(languageMap[k])
			assert.NoError(t, err)

			analyzeURL(t, url, o, total)
		}

		o := Options{language: allLanguages}

		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)

		analyzeURL(t, url, o, 231) // 231 = 21!
	})
}

// this function is really to just reduce the repetitive blocks of code
// within the tests for conditions and languages. These mappings relate strings
// to integers. With that in mind, we can check for appropriate mapping
// my passing a total, subtracting the passed values, and expecting a 0 result.
func analyzeURL(t *testing.T, url string, o Options, total int) {
	t.Helper()

	beginQ := strings.Index(url, "?")
	urlArgs := strings.Split(url[beginQ+1:], "&")
	for _, arg := range urlArgs {
		pieces := strings.Split(arg, "=")
		key := pieces[0]
		if key != "query" && key != "sort" {
			val, err := strconv.Atoi(pieces[1])
			assert.NoError(t, err)
			total -= val
		}
	}
	assert.Equal(t, 0, total)
}