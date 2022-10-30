package main

import "github.com/goyalmunish/reminder/cmd/reminder"

var version string

func main() {
	// go utils.Spinner(100 * time.Millisecond)
	reminder.Flow()
}
