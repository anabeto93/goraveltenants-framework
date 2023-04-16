package events

import (
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
	"github.com/anabeto93/goraveltenants/contracts"
)

var _ eventcontracts.TenancyEvent = &BootstrappingTenancy{}

type BootstrappingTenancy struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *BootstrappingTenancy) Name() string {
	return "BootstrappingTenancy"
}

func NewBootstrappingTenancyEvent(tenant contracts.Tenant) *BootstrappingTenancy {
	return eventcontracts.NewTenancyEvent("BootstrappingTenancy", tenant).(*BootstrappingTenancy)
}
