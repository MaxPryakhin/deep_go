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
	errors []error
}

func (e *MultiError) Error() string {
	l := len(e.errors)
	if l == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occured:\n", l))
	for _, err := range e.errors {
		sb.WriteString("\t* " + err.Error())
	}
	sb.WriteString("\n")
	return sb.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError
	if err != nil {
		if me, ok := err.(*MultiError); ok {
			multiErr = me
		} else {
			multiErr = &MultiError{errors: []error{err}}
		}
	} else {
		multiErr = &MultiError{}
	}
	multiErr.errors = append(multiErr.errors, errs...)
	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
