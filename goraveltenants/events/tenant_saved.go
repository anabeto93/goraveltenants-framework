package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &TenantSaved{}

type TenantSaved struct {
	eventcontracts.BaseTenantEvent
}

func NewTenantSaved(tenant contracts.Tenant) *TenantSaved {
	return eventcontracts.NewTenantEvent("TenantSaved", tenant).(*TenantSaved)
}
