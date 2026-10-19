package errors

import "errors"

var (
	ErrNotFound    = errors.New("metric not found")
	ErrNotInitDB   = errors.New("not initialized db")
	ErrSendMetrics = errors.New("error agent send metric")
	ErrNotFIleName = errors.New("error not fIle name")
	ErrNotUrl      = errors.New("error not url")
)
