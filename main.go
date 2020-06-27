package main

import (
	"context"
	"fmt"
)

func main() {
	url := "https://newyork.craigslist.org/search/pta?srchType=T"

	c := NewClient("newyork")

	result, err := c.GetListings(context.TODO(), url)
	if err != nil {
		panic(err)
	}

	fmt.Println("the reuslt:", result)

	for !result.Done {
		fmt.Println("the first result:", result.Listings[0].Title)
		fmt.Println("the last result:", result.Listings[len(result.Listings)-1].Title)

		result, err = result.Next(context.TODO())
		if err != nil {
			panic(err)
		}
	}

}
