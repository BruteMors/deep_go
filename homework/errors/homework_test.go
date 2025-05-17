package main

import (
	"errors"
	"fmt"
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

	flatError := fmt.Sprintf("%d errors occured:\n", len(e.Errors))
	for _, err := range e.Errors {
		flatError += fmt.Sprintf("\t* %s", err.Error())
	}
	flatError += "\n"

	return flatError
}

func Append(err error, errs ...error) *MultiError {
	filteredErrs := make([]error, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			filteredErrs = append(filteredErrs, e)
		}
	}

	if err == nil && len(filteredErrs) == 0 {
		return nil
	}

	var multiErr *MultiError
	if errors.As(err, &multiErr) {
		multiErr.Errors = append(multiErr.Errors, filteredErrs...)
		return multiErr
	}

	resultErrs := make([]error, 0, len(filteredErrs)+1)

	if err != nil {
		resultErrs = append(resultErrs, err)
	}

	resultErrs = append(resultErrs, filteredErrs...)

	return &MultiError{
		Errors: resultErrs,
	}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
