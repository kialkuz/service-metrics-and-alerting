package errors

import "errors"

var (
	ErrNotFound  = errors.New("metric not found")
	ErrNotInitDB = errors.New("not initialized db")
)
