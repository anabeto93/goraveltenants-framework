package exceptions

import (
    "fmt"
    "github.com/anabeto93/goraveltenants/contracts"
)

type TenantCouldNotBeIdentifiedById struct {
    tenantID string
}

func NewTenantCouldNotBeIdentifiedById(tenantID string) *TenantCouldNotBeIdentifiedById {
    return &TenantCouldNotBeIdentifiedById{tenantID: tenantID}
}

func (e *TenantCouldNotBeIdentifiedById) Error() string {
    return fmt.Sprintf("Tenant could not be identified with tenant_id: %s", e.tenantID)
}

func (e *TenantCouldNotBeIdentifiedById) Unwrap() error {
    return contracts.NewTenantCouldNotBeIdentifiedException(e.Error())
}
