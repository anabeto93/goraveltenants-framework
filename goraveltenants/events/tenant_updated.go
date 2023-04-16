package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &TenantUpdated{}

type TenantUpdated struct {
	eventcontracts.BaseTenantEvent
}

func NewTenantUpdated(tenant contracts.Tenant) *TenantUpdated {
	return eventcontracts.NewTenantEvent("TenantUpdated", tenant).(*TenantUpdated)
}
