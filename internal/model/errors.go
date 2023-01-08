package model

import "errors"

var (
	ErrorConflictFile = errors.New("created _CONFLICT file")
	ErrorMutexLockOn  = errors.New("Mutex Lock is ON; there is already a session running!")
)
