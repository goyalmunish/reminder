package appinfo

import (
	"os"
	"path"
)

type Options struct {
	DataFile string
}

func DefaultOptions() *Options {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		homePath = "~"
	}
	dataFilePath := path.Join(homePath, "reminder", "data.json")
	return &Options{
		DataFile: dataFilePath,
	}
}
