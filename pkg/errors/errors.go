package errors

import "errors"

var (
	ErrNotFound    = errors.New("metric not found")
	ErrNotInitDB   = errors.New("not initialized db")
	ErrNotFIleName = errors.New("error not fIle name")
	ErrNotURL      = errors.New("error not url")
	ErrAppStopped  = errors.New("application stopped")
)
