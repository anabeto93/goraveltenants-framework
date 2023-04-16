package contracts

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/contracts/event"
)

type TenantEvent interface {
	GetTenant() contracts.Tenant
	Name() string
	Handle(args []event.Arg) ([]event.Arg, error)
}

type BaseTenantEvent struct {
	tenant contracts.Tenant
	name   string
}

func (e *BaseTenantEvent) GetTenant() contracts.Tenant {
	return e.tenant
}

func (e *BaseTenantEvent) Name() string {
	return e.name
}

func (e *BaseTenantEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

func NewTenantEvent(name string, tenant contracts.Tenant) TenantEvent {
	return &BaseTenantEvent{
		tenant: tenant,
		name:   name,
	}
}
