package contracts

import "github.com/anabeto93/goraveltenants/contracts"

type TenancyEvent interface {
	Name() string
	GetTenant() contracts.Tenant
}

type BaseTenancyEvent struct {
	name string
	tenant contracts.Tenant
}

func (e *BaseTenancyEvent) Name() string {
	return e.name
}

func (e *BaseTenancyEvent) GetTenant() contracts.Tenant {
	return e.tenant
}

func NewBaseTenancy(name string, tenant contracts.Tenant) *BaseTenancyEvent {
	return &BaseTenancyEvent{name, tenant}
}

func NewTenancyEvent(name string, tenant contracts.Tenant) TenancyEvent {
	return NewBaseTenancy(name, tenant)
}
