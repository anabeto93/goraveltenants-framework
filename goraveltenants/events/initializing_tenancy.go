package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &InitializingTenancy{}

type InitializingTenancy struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *InitializingTenancy) Name() string {
	return "InitializingTenancy"
}

func NewInitializingTenancyEvent(tenant contracts.Tenant) *InitializingTenancy {
	return eventcontracts.NewTenancyEvent("InitializingTenancy", tenant).(*InitializingTenancy)
}
