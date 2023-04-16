package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DeletingTenant{}

type DeletingTenant struct {
	eventcontracts.BaseTenantEvent
}

func NewDeletingTenant(tenant contracts.Tenant) *DeletingTenant {
	return eventcontracts.NewTenantEvent("DeletingTenant", tenant).(*DeletingTenant)
}
