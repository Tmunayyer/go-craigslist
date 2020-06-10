package main

import (
	"context"
	"fmt"
)

func noErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	c, err := NewClient(context.TODO(), "newyork")
	noErr(err)

	url := "https://newyork.craigslist.org/search/sss?query=xbox&sort=rel&srchType=T&postedToday=1"

	result, err := c.Search(context.TODO(), url)
	noErr(err)

	fmt.Println("the result:", result)

}
