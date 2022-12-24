package calendar

import (
	"os"
	"path"
)

type Options struct {
	CredentialFile string
}

func DefaultOptions() *Options {
	homePath := os.Getenv("HOME")
	if homePath != "" {
		homePath = "~"
	}
	return &Options{
		CredentialFile: path.Join(homePath, "credentials.json"),
	}
}
