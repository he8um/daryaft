package downloader

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var ErrCancelled = errors.New("download cancelled")

const (
	defaultRetryBaseDelay = time.Second
	defaultRetryMaxDelay  = 8 * time.Second
)

type Sleeper func(time.Duration)

type RetryPolicy struct {
	Retries int
	Sleep   Sleeper
}

func (p RetryPolicy) MaxAttempts() int {
	return maxAttempts(p.Retries)
}

type httpStatusError struct {
	StatusCode int
	Status     string
}

type nonRetryableError struct {
	err error
}

func (e nonRetryableError) Error() string {
	return e.err.Error()
}

func (e nonRetryableError) Unwrap() error {
	return e.err
}

func (e httpStatusError) Error() string {
	if isRetryableStatus(e.StatusCode) {
		return "temporary server error: " + e.Status
	}
	return "server returned " + e.Status
}

func maxAttempts(retries int) int {
	if retries < 0 {
		retries = 0
	}
	return retries + 1
}

func BackoffDelay(failedAttempt int) time.Duration {
	if failedAttempt <= 0 {
		return 0
	}

	delay := defaultRetryBaseDelay
	for attempt := 1; attempt < failedAttempt; attempt++ {
		delay *= 2
		if delay >= defaultRetryMaxDelay {
			return defaultRetryMaxDelay
		}
	}
	return delay
}

func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrCancelled) || errors.Is(err, context.Canceled) {
		return false
	}

	var permanent nonRetryableError
	if errors.As(err, &permanent) {
		return false
	}

	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		return isRetryableStatus(statusErr.StatusCode)
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	return os.IsTimeout(err)
}

func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
