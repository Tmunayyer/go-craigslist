package main

func main() {
	c := Client{}
	_, err := c.ListCategories()
	noErr(err)
	_, err = c.ListLocations()
	noErr(err)

	c.PrintCategories()
}
