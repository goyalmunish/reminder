package models

import (
	"os"
	"path"
)

var DataFile = path.Join(os.Getenv("HOME"), "reminder.json")
