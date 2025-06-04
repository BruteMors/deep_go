package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.Errors)))

	for _, err := range e.Errors {
		sb.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}

	sb.WriteString("\n")

	return sb.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError

	if !errors.As(err, &multiErr) && err != nil {
		multiErr = &MultiError{Errors: []error{err}}
	}

	for _, e := range errs {
		if e != nil {
			if multiErr == nil {
				multiErr = &MultiError{}
			}
			multiErr.Errors = append(multiErr.Errors, e)
		}
	}

	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
