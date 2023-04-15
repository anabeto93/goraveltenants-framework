package exceptions

import "fmt"

type NoConnectionSetException struct {
	manager string
}

func NewNoConnectionSetException(manager string) NoConnectionSetException {
	return NoConnectionSetException{
		manager: manager,
	}
}

func (e NoConnectionSetException) Error() string {
	return fmt.Sprintf("No connection was set on this %s instance.", e.manager)
}
