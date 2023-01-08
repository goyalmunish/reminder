package main

import (
	"github.com/goyalmunish/reminder/cmd/reminder"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func main() {
	// go utils.Spinner(100 * time.Millisecond)
	err := reminder.Run()
	utils.LogError(err)
}
