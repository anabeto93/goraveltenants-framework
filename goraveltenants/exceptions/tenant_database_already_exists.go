package exceptions

import (
	"fmt"
	"github.com/anabeto93/goraveltenants/contracts"
)

type TenantDatabaseAlreadyExistsException struct {
	database string
}

func NewTenantDatabaseAlreadyExistsException(database string) *TenantDatabaseAlreadyExistsException {
	return &TenantDatabaseAlreadyExistsException{
		database: database,
	}
}

func (e *TenantDatabaseAlreadyExistsException) Error() string {
	return fmt.Sprintf("Tenant database %s already exists", e.database)
}

func (e *TenantDatabaseAlreadyExistsException) Reason() string {
	return fmt.Sprintf("Database %s already exists.", e.database)
}

func (e *TenantDatabaseAlreadyExistsException) Unwrap() error {
	return contracts.NewTenantCannotBeCreatedException(e.Error())
}
