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

	c, err := NewClient(context.TODO())
	noErr(err)

	c.PrintCategories()

	q, err := c.BuildQuery("nyc", "ela", "ps4", Filters{})
	noErr(err)

	result, err := c.Search(context.TODO(), q)
	noErr(err)

	fmt.Println("the result:", result)

}
