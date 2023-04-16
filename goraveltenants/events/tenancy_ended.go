package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &TenancyEnded{}

type TenancyEnded struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *TenancyEnded) Name() string {
	return "TenancyEnded"
}

func NewTenancyEndedEvent(tenant contracts.Tenant) *TenancyEnded {
	return eventcontracts.NewTenancyEvent("TenancyEnded", tenant).(*TenancyEnded)
}
