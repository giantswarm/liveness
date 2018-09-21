package main

import (
	"net/http"

	"github.com/giantswarm/microerror"
)

var serverClosedError = microerror.New("server closed")

// IsServerClosed asserts serverClosedError.
func IsServerClosed(err error) bool {
	c := microerror.Cause(err)

	if c == http.ErrServerClosed {
		return true
	}

	if c == serverClosedError {
		return true
	}

	return false
}
