package exceptions

import (
    "fmt"
    "github.com/anabeto93/goraveltenants/contracts"
)

type TenantCouldNotBeIdentifiedByPathException struct {
    tenantID string
}

func NewTenantCouldNotBeIdentifiedByPathException(tenantID string) *TenantCouldNotBeIdentifiedByPathException {
    return &TenantCouldNotBeIdentifiedByPathException{
		tenantID: tenantID,
    }
}

func (e *TenantCouldNotBeIdentifiedByPathException) Error() string {
    return fmt.Sprintf("Tenant could not be identified on path with tenant_id: %s", e.tenantID)
}

func (e *TenantCouldNotBeIdentifiedByPathException) Unwrap() error {
    return contracts.NewTenantCouldNotBeIdentifiedException(e.Error())
}
