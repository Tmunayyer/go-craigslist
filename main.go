package main

func main() {
	c := Client{}
	_, err := c.ListCategories()
	noErr(err)
	locations, err := c.ListLocations()
	noErr(err)
}
