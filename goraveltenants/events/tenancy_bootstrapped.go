package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &TenancyBootstrapped{}

type TenancyBootstrapped struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *TenancyBootstrapped) Name() string {
	return "TenancyBootstrapped"
}

func NewTenancyBootstrappedEvent(tenant contracts.Tenant) *TenancyBootstrapped {
	return eventcontracts.NewTenancyEvent("TenancyBootstrapped", tenant).(*TenancyBootstrapped)
}
