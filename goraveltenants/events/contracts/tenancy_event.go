package contracts

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/contracts/event"
)

type TenancyEvent interface {
	Name() string
	GetTenant() contracts.Tenancy
	Handle(args []event.Arg) ([]event.Arg, error)
}

type BaseTenancyEvent struct {
	name   string
	tenant contracts.Tenancy
}

func (e *BaseTenancyEvent) Name() string {
	return e.name
}

func (e *BaseTenancyEvent) GetTenant() contracts.Tenancy {
	return e.tenant
}

func (e *BaseTenancyEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

func NewBaseTenancy(name string, tenant contracts.Tenancy) *BaseTenancyEvent {
	return &BaseTenancyEvent{name, tenant}
}

func NewTenancyEvent(name string, tenant contracts.Tenancy) TenancyEvent {
	return NewBaseTenancy(name, tenant)
}
