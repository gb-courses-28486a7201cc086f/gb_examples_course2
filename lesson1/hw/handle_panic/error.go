package main

import (
	"fmt"
	"time"
)

type CustomError struct {
	raiseTime time.Time
	message   string
}

func (e CustomError) Error() string {
	formatTS := e.raiseTime.Format(time.RFC3339)
	msg := fmt.Sprintf("%s ERROR: %s", formatTS, e.message)
	return msg
}

func NewError(message string) error {
	err := CustomError{time.Now(), message}
	return err
}
