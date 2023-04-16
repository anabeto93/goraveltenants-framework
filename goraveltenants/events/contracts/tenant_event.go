package contracts

import "github.com/anabeto93/goraveltenants/contracts"

type TenantEvent interface {
    GetTenant() contracts.Tenant
    Name() string
}

type BaseTenantEvent struct {
    tenant contracts.Tenant
    name string
}

func (e *BaseTenantEvent) GetTenant() contracts.Tenant {
    return e.tenant
}

func (e *BaseTenantEvent) Name() string {
    return e.name
}

func NewTenantEvent(name string, tenant contracts.Tenant) TenantEvent {
    return &BaseTenantEvent{
        tenant: tenant,
        name: name,
    }
}
