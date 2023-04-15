package exceptions

import (
	"fmt"

	"github.com/anabeto93/goraveltenants/contracts"
)

type TenantDatabaseUserAlreadyExistsException struct {
	user string
}

func NewTenantDatabaseUserAlreadyExistsException(user string) *TenantDatabaseUserAlreadyExistsException {
	return &TenantDatabaseUserAlreadyExistsException{user: user}
}

func (e *TenantDatabaseUserAlreadyExistsException) Error() string {
	return fmt.Sprintf("Database user %s already exists", e.user)
}

func (e *TenantDatabaseUserAlreadyExistsException) Reason() string {
	return fmt.Sprintf("Database user %s already exists", e.user)
}

func (e *TenantDatabaseUserAlreadyExistsException) Unwrap() error {
	return contracts.NewTenantCannotBeCreatedException(e.Error())
}
