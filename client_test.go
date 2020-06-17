package gocraigslist

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatURL(t *testing.T) {
	client := NewClient("newyork")

	t.Run("no options provided, basic term", func(t *testing.T) {
		expected := "https://newyork.craigslist.org/search/sss?query=xbox&sort=rel"
		url := client.FormatURL("xbox", Options{})
		assert.Equal(t, expected, url)
	})

	t.Run("properly escaping term", func(t *testing.T) {
		expected := "https://newyork.craigslist.org/search/sss?query=xbox+123+.+%23%24%25&sort=rel"
		url := client.FormatURL("xbox 123 . #$%", Options{})
		assert.Equal(t, expected, url)
	})

	t.Run("should overwright location if provided by options", func(t *testing.T) {
		o := Options{Location: "testing"}
		expected := "https://testing.craigslist.org/search/sss?query=xbox&sort=rel"
		url := client.FormatURL("xbox", o)
		assert.Equal(t, expected, url)
	})

	t.Run("should overwright category if provided by options", func(t *testing.T) {
		o := Options{Category: "aaa"}
		expected := "https://newyork.craigslist.org/search/aaa?query=xbox&sort=rel"

		url := client.FormatURL("xbox", o)
		assert.Equal(t, expected, url)
	})

	t.Run("should account for bool options", func(t *testing.T) {
		o := Options{
			SrchType:          true,
			HasPic:            true,
			PostedToday:       true,
			BundleDuplicates:  true,
			CryptoCurrencyOK:  true,
			DeliveryAvailable: true,
		}

		url := client.FormatURL("xbox", o)

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
		if o.SrchType {
			assert.Equal(t, "T", argMap["srchType"])
			numberOfOptions++
		}

		if o.HasPic {
			assert.Equal(t, "1", argMap["hasPic"])
			numberOfOptions++
		}

		if o.PostedToday {
			assert.Equal(t, "1", argMap["postedToday"])
			numberOfOptions++
		}

		if o.PostedToday {
			assert.Equal(t, "1", argMap["postedToday"])
			numberOfOptions++
		}

		if o.BundleDuplicates {
			assert.Equal(t, "1", argMap["bundleDuplicates"])
			numberOfOptions++
		}

		if o.CryptoCurrencyOK {
			assert.Equal(t, "1", argMap["crypto_currency_ok"])
			numberOfOptions++
		}

		if o.DeliveryAvailable {
			assert.Equal(t, "1", argMap["delivery_available"])
			numberOfOptions++
		}

		// the query and rel arg are passed always, minus 2 from length of argMap
		assert.Equal(t, numberOfOptions, len(argMap)-1)
	})

	t.Run("accounts for min and max price", func(t *testing.T) {
		o := Options{
			MinPrice: "100",
			MaxPrice: "500",
		}

		url := client.FormatURL("xbox", o)

		beginQ := strings.Index(url, "?")
		urlArgs := strings.Split(url[beginQ+1:], "&")
		argMap := map[string]string{}
		for _, arg := range urlArgs {
			pieces := strings.Split(arg, "=")
			argMap[pieces[0]] = pieces[1]
		}

		assert.Equal(t, o.MinPrice, argMap["min_price"])
		assert.Equal(t, o.MaxPrice, argMap["max_price"])
	})

	t.Run("accounts for conditions", func(t *testing.T) {
		for _, test := range []struct {
			given    Options
			expected int
		}{
			{
				given:    Options{Condition: []string{"new"}},
				expected: 10,
			},
			{
				given:    Options{Condition: []string{"like new"}},
				expected: 20,
			},
			{
				given:    Options{Condition: []string{"excellent"}},
				expected: 30,
			},
			{
				given:    Options{Condition: []string{"good"}},
				expected: 40,
			},
			{
				given:    Options{Condition: []string{"fair"}},
				expected: 50,
			},
			{
				given:    Options{Condition: []string{"salvage"}},
				expected: 60,
			},
			{
				given:    Options{Condition: []string{"new", "like new", "excellent", "good", "fair", "salvage"}},
				expected: 210, // 10 + 20 + 30 + 40 + 50 + 60
			},
		} {
			url := client.FormatURL("xbox", test.given)
			analyzeURL(t, url, test.given, test.expected)
		}
	})

	t.Run("accounts for languages", func(t *testing.T) {
		languageMap := map[string]string{"af": "1", "ca": "2", "da": "3", "de": "4", "en": "5", "es": "6", "fi": "7", "fr": "8", "it": "9", "nl": "10", "no": "11", "pt": "12", "sv": "13", "tl": "14", "tr": "15", "zh": "16", "ar": "17", "ja": "18", "ko": "19", "ru": "20", "vi": "21"}

		allLanguages := []string{}
		for k := range languageMap {
			// push to a slice to use for testing all langs later
			allLanguages = append(allLanguages, k)

			o := Options{Language: []string{k}}

			url := client.FormatURL("xbox", o)

			total, err := strconv.Atoi(languageMap[k])
			assert.NoError(t, err)

			analyzeURL(t, url, o, total)
		}

		o := Options{Language: allLanguages}

		url := client.FormatURL("xbox", o)

		analyzeURL(t, url, o, 231) // 231 = 21!
	})
}

// analyzeURL is specifically for conditions and languages tests. Lots of repeated code
// could be moved into here. The core idea is to
// 		1. break up the URL
// 		2. pull out the arguments
// 		3. find out what FormatURL set them to
//		4. compare that to what it should be
// A lot of this logic is done through simple math since the craigslist format
// for most of these setting is to use integers.
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
