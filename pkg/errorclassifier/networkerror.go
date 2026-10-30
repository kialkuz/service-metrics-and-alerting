package errorclassifier

import (
	"context"
	"errors"
	"net"
	"net/http"
)

type HTTPError struct {
	error
	StatusCode int
	Status     string
}

type NetworkErrorClassifier struct{}

func NewNetworkErrorClassifier() ErrorClassifier {
	return &NetworkErrorClassifier{}
}

func (c NetworkErrorClassifier) Classify(err error) ErrorClassification {
	if errors.Is(err, context.Canceled) {
		return NonRetriable
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return Retriable
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return Retriable
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return Retriable
	}

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.StatusCode {
		case http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusNotImplemented:
			return NonRetriable

		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return Retriable
		}
	}

	return NonRetriable
}
