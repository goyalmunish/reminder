package utils

import "io"

func PrivateAskBoolean(msg string, in io.Reader) (bool, error) {
	return askBoolean(msg, in)
}
