package container

import (
	"errors"
)

type CircularDependencyError struct {
	err error
}

func (e CircularDependencyError) Error() string {
	return e.err.Error()
}

func (e CircularDependencyError) Unwrap() error {
	return e.err
}

func NewCircularDependencyError(msg string) error {
	return CircularDependencyError{err: errors.New(msg)}
}
