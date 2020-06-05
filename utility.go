package main


func noErr(err error) {
	if err != nil {
		panic(err)
	}
}