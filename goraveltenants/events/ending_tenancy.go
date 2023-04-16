package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &EndingTenancy{}

type EndingTenancy struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *EndingTenancy) Name() string {
	return "EndingTenancy"
}

func NewEndingTenancyEvent(tenant contracts.Tenancy) *EndingTenancy {
	return eventcontracts.NewTenancyEvent("EndingTenancy", tenant).(*EndingTenancy)
}
