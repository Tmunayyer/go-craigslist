# go-craigslist
This is a small library to programatically search craigslist.com. Inspired by [this node library](https://github.com/brozeph/node-craigslist).

## Table of Contents

- [Quickstart](#quickstart)
- [Options](#options)
- [Categories and Locations](#categoriesandlocations)

## Quickstart
```go
package main

import (
	"context"
	"fmt"

	"github.com/tmunayyer/gocraigslist"
)

func main() {
    // create a new client, pass in a location
	client := gocraigslist.NewClient("newyork")

	// use your own url from a search on craigslist.com
	listings, err := client.GetListings(context.TODO(), "https://newyork.craigslist.org/d/antiques/search/ata")
	if err != nil {
		panic(err)
	}
	fmt.Println("the antique listings:", listings)
}
```

## Options
| propert name          | type      | required | default      | description |
|-----------------------|-----------|----------|--------------|-------------|
|  location             | string    | false    | *init value  | defaults to location provided on intialization, providing location here will overrides init value |
|  category             | string    | false    | all          | all, owner, dealer **attention**: incompatible for specific categories |
|  hasPic               | bool      | false    | false        | true or false |
|  postedToday          | bool      | false    | false        | true or false |
|  bundleDuplicates     | bool      | false    | false        | true or false |
|  cryptoCurrencyOK     | bool      | false    | false        | true or false |
|  deliveryAvailable    | bool      | false    | false        | true or false |
|  minPrice             | string    | false    | 0            | example: "100" |
|  maxPrice             | string    | false    | 0            | example: "500" |
|  lanaguage            | []string  | false    | []string     | new, like new, excellent, good, fair, salvage |
|  condition            | []string  | false    | []string     | af, ca, da, de, en, es, fi, fr, it, nl, no, pt, sv, tl, tr, zh, ar, ja, ko, ru, vi |

## Categories and Locations

Resource: https://www.craigslist.org/about/reference

The above resource details an api to discover categories and locations. This api is a bit unreliable/inconsistent. The most precise way to find a specfic category to search is to go directly to craigslist homepage and select it. You can then find the three letter code in the url. See example below.

### Example:
```
https://sfbay.craigslist.org/d/auto-parts/search/pta
```
The location is "sfbay". The category code is "pta".

