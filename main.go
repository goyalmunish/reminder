package main

import (
	"github.com/goyalmunish/reminder/cmd/reminder"
	"github.com/goyalmunish/reminder/pkg/utils"
)

func main() {
	err := reminder.Run()
	utils.LogError(err)
}
