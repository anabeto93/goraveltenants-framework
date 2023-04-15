package contracts

import "github.com/anabeto93/goraveltenants/contracts"

type TenantEvent interface {
}

type tenantEvent struct {
    Tenant contracts.Tenant
}

func (e *tenantEvent) GetTenant() contracts.Tenant {
    return e.Tenant
}

func NewTenantEvent(tenant contracts.Tenant) TenantEvent {
    return &tenantEvent{Tenant: tenant}
}
