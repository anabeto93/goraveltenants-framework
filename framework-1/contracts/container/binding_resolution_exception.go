package container

import (
	"errors"
)

type BindingResolutionError struct {
	err error
}

func (e BindingResolutionError) Error() string {
	return e.err.Error()
}

func (e BindingResolutionError) Unwrap() error {
	return e.err
}

func NewBindingResolutionError(msg string) error {
	return BindingResolutionError{err: errors.New(msg)}
}
