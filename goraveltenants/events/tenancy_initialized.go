package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &TenancyInitialized{}

type TenancyInitialized struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *TenancyInitialized) Name() string {
	return "TenancyInitialized"
}

func NewTenancyInitializedEvent(tenant contracts.Tenant) *TenancyInitialized {
	return eventcontracts.NewTenancyEvent("TenancyInitialized", tenant).(*TenancyInitialized)
}
