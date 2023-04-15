package exceptions

import "fmt"

type DatabaseManagerNotRegisteredException struct {
	driver string
}

func NewDatabaseManagerNotRegisteredException(driver string) error {
	return &DatabaseManagerNotRegisteredException{driver: driver}
}

func (e *DatabaseManagerNotRegisteredException) Error() string {
	return fmt.Sprintf("Database manager for driver %s is not registered.", e.driver)
}
