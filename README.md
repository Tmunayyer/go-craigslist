# go-craigslist
A Go API to query craigslist.

## To Use

```go
package main

import (
	"context"
	"fmt"

	"github.com/tmunayyer/gocraigslist"
)

func main() {
	client := gocraigslist.NewClient("newyork")

	// with a url
	listings, err := client.GetListings(context.TODO(), "https://newyork.craigslist.org/d/antiques/search/ata")
	if err != nil {
		panic(err)
	}
	fmt.Println("the antique listings:", listings)

	// or create your own url
	url := client.FormatURL("comfy couch", gocraigslist.Options{
		location:  "sfbay",
		category:  "fua", // furniture
		maxPrice:  "500",
		condition: []string{"new", "like new"},
	})

	listings, err = client.GetListings(context.TODO(), url)
	if err != nil {
		panic(err)
	}
	fmt.Println("the couch listings in SF:", listings)
}
```

## Options
Optionsal filters with tuple values are represented as [input value, mapped value]

| propert name          | type      | required | default      | description |
|-----------------------|-----------|----------|--------------|-------------|
|  location             | string    | false    | *init value  | defaults to location provided on intialization, providing location here will overrides init value |
|  category             | string    | false    | "sss"        | [all, sss], [owner, sso], [dealer, ssq] **attention**: this only works for default search (sss), not specific categories |
|  hasPic               | bool      | false    | false        | true or false |
|  postedToday          | bool      | false    | false        | true or false |
|  bundleDuplicates     | bool      | false    | false        | true or false |
|  cryptoCurrencyOK     | bool      | false    | false        | true or false |
|  deliveryAvailable    | bool      | false    | false        | true or false |
|  minPrice             | string    | false    | 0            | example: "100" |
|  maxPrice             | string    | false    | 0            | example: "500" |
|  lanaguage            | []string  | false    | []string     | [new, 10], [like new, 20], [excellent, 30], [good, 40], [fair, 50], [salvage, 60] |
|  condition            | []string  | false    | []string     | [af, 1], [ca, 2], [da, 3], [de, 4], [en, 5], [es, 6], [fi, 7], [fr, 8], [it, 9], [nl, 10], [no, 11], [pt, 12], [sv, 13], [tl, 14], [tr, 15], [zh, 16], [ar, 17], [ja, 18], [ko, 19], [ru, 20], [vi, 21] |

## Categories and Locations

Resource: https://www.craigslist.org/about/reference

The above resource details an api to discover categories and locations. This api is a bit unreliable/inconsistent. The most precise way to find a specfic category to search is to go directly to craigslist homepage and select it. You can then find the three letter code in the url. See example below.

### Example:
```
https://sfbay.craigslist.org/d/auto-parts/search/pta
```
The Category code here is "pta".