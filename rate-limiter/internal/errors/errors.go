package errors

import "errors"

var (
	ErrTooManyRequests = errors.New("you have reached the maximum number of requests or actions allowed within a certain time frame")
)
