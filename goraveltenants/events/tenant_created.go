package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &TenantCreated{}

type TenantCreated struct {
	eventcontracts.BaseTenantEvent
}

func NewTenantCreated(tenant contracts.Tenant) *TenantCreated {
	return eventcontracts.NewTenantEvent("TenantCreated", tenant).(*TenantCreated)
}
