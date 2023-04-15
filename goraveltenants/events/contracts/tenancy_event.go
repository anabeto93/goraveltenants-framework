package contracts

import "github.com/anabeto93/goraveltenants/contracts"

type TenancyEvent interface {
	Name() string
	GetTenant() contracts.Tenant
}

type tenancyEvent struct {
	name string
	tenant contracts.Tenant
}

func (e *tenancyEvent) Name() string {
	return e.name
}

func (e *tenancyEvent) GetTenant() contracts.Tenant {
	return e.tenant
}

func NewTenancyEvent(name string, tenant contracts.Tenant) TenancyEvent {
	return &tenancyEvent{name, tenant}
}
