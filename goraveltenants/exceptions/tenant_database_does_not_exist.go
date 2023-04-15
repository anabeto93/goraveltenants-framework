package exceptions

import "fmt"

type TenantDatabaseDoesNotExistException struct {
    database string
}

func NewTenantDatabaseDoesNotExistException(database string) *TenantDatabaseDoesNotExistException {
    return &TenantDatabaseDoesNotExistException{database: database}
}

func (e *TenantDatabaseDoesNotExistException) Error() string {
    return fmt.Sprintf("Database %s does not exist.", e.database)
}
