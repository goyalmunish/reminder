package model

import "errors"

var (
	ErrorConflictFile              = errors.New("Created _CONFLICT file")
	ErrorMutexLockOn               = errors.New("Mutex Lock is ON; there is already a session running!")
	ErrorInteractiveProcessSkipped = errors.New("Skipped running the interactive process. Try again!")
)
