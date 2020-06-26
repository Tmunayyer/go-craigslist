package main

import "context"

func main() {
	url := "https://newyork.craigslist.org/search/pta?srchType=T"

	c := NewClient("newyork")

	c.GetMultipageListings(context.TODO(), url)
}
