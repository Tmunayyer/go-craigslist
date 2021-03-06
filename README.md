# go-craigslist
This is a small library to programatically search craigslist.com. Inspired by [this node library](https://github.com/brozeph/node-craigslist).

## Table of Contents

- [Quickstart](#quickstart)
- [Options](#options)
- [Categories and Locations](#categoriesandlocations)
- [Documentation](https://godoc.org/github.com/Tmunayyer/go-craigslist)

## Quickstart
```go
package main

import (
	"context"
	"fmt"

	"github.com/tmunayyer/gocraigslist"
)

func main() {
	client := gocraigslist.NewClient("newyork")

	result, err := client.GetListings(context.TODO(), "https://newyork.craigslist.org/d/antiques/search/ata")
	if err != nil {
		panic(err)
	}
	fmt.Println("the antique listings:", result.Listings)
}
```

## Options
| propert name          | type      | required | default      | description |
|-----------------------|-----------|----------|--------------|-------------|
|  location             | string    | false    | *init value  | defaults to location provided on intialization, providing location here will overrides init value |
|  category             | string    | false    | all          | *see section Categories and Locations |
|  srchType             | string    | false    | "all"        | all, owner, dealer **attention**: incompatible for specific categories |
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

The above resource details an api to discover categories and locations. This api is a bit unreliable/inconsistent. The most precise way to find a specfic category is to go directly to craigslist's homepage and select from the index. You can then find the three letter code in the url. See example below.

### Example:
```
https://sfbay.craigslist.org/d/auto-parts/search/pta
```
The location is "sfbay". The category code is "pta".

