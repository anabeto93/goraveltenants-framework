package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &BootstrappingTenancy{}

type BootstrappingTenancy struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *BootstrappingTenancy) Name() string {
	return "BootstrappingTenancy"
}

func NewBootstrappingTenancyEvent(tenant contracts.Tenancy) *BootstrappingTenancy {
	return eventcontracts.NewTenancyEvent("BootstrappingTenancy", tenant).(*BootstrappingTenancy)
}
