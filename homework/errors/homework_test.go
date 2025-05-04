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
	errs []error
}

func (e *MultiError) Error() string {
	switch len(e.errs) {
	case 0:
		return ""
	case 1:
		return e.errs[0].Error()
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errs)))

	for _, err := range e.errs {
		sb.WriteString(fmt.Sprintf("\t* %v", err))
	}

	sb.WriteByte('\n')
	return sb.String()
}

func (e *MultiError) Append(errs ...error) {
	for _, v := range errs {
		if v == nil {
			continue
		}
		if m, ok := v.(*MultiError); ok {
			e.errs = append(e.errs, m.errs...)
		} else {
			e.errs = append(e.errs, v)
		}
	}
}

func (e *MultiError) Unwrap() []error {
	return e.errs
}

func Append(err error, errs ...error) *MultiError {
	var res MultiError
	res.Append(err)
	res.Append(errs...)
	return &res
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	err3 := errors.New("error 3")
	err = Append(err, err3)
	assert.ErrorIs(t, err, err3)

	var merr *MultiError
	assert.ErrorAs(t, err, &merr)
	assert.Contains(t, merr.errs, err3)
}
