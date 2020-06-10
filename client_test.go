package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomLanguage(n int) []string {
	langBank := []string{"af", "ca", "da", "de", "en", "es", "fi", "fr", "it", "nl", "no", "pt", "sv", "tl", "tr", "zh", "ar", "ja", "ko", "ru", "vi"}

	out := []string{}
	for i := 0; i < n; i++ {
		index := rand.Intn(len(langBank))
		selected := langBank[index]

		langBank = append(langBank[0:index], langBank[index+1:len(langBank)]...)
		out = append(out, selected)
	}

	return out
}

func randomCondition(n int) []string {
	langBank := []string{"new", "like new", "excellent", "good", "fair", "salvage"}

	out := []string{}
	for i := 0; i < n; i++ {
		index := rand.Intn(len(langBank))
		selected := langBank[index]

		langBank = append(langBank[0:index], langBank[index+1:len(langBank)]...)
		out = append(out, selected)
	}

	return out
}

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
		url, err := client.FormatURL("xbox", Options{})
		assert.NoError(t, err)
		expected := "https://newyork.craigslist.org/search/sss?query=xbox&sort=rel"

		assert.Equal(t, expected, url)
	})

	t.Run("properly escaping term", func(t *testing.T) {
		url, err := client.FormatURL("xbox 123 . #$%", Options{})
		assert.NoError(t, err)
		expected := "https://newyork.craigslist.org/search/sss?query=xbox+123+.+%23%24%25&sort=rel"

		assert.Equal(t, expected, url)
	})

	t.Run("should overwright location if provided by options", func(t *testing.T) {
		o := Options{
			location: "testing",
		}
		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)

		expected := "https://testing.craigslist.org/search/sss?query=xbox&sort=rel"

		assert.Equal(t, expected, url)
	})

	t.Run("should overwright category if provided by options", func(t *testing.T) {
		o := Options{
			category: "aaa",
		}
		url, err := client.FormatURL("xbox", o)
		assert.NoError(t, err)

		expected := "https://newyork.craigslist.org/search/aaa?query=xbox&sort=rel"

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

		beginQ := strings.Index(url, "?")
		urlArgs := strings.Split(url[beginQ+1:], "&")
		fmt.Println("the url args:", urlArgs)
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
			fmt.Println("cryptoOK", argMap["crypto_currency_ok"])
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
		o := Options{
			condition: []string{"new", "like new", "excellent", "good", "fair", "salvage"},
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
}

func index(slice []string, target string) int {
	for i, val := range slice {
		if val == target {
			return i
		}
	}
	return -1
}
