package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &TenantDeleted{}

type TenantDeleted struct {
	eventcontracts.BaseTenantEvent
}

func NewTenantDeleted(tenant contracts.Tenant) *TenantDeleted {
	return eventcontracts.NewTenantEvent("TenantDeleted", tenant).(*TenantDeleted)
}
