package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &UpdatingTenant{}

type UpdatingTenant struct {
	eventcontracts.BaseTenantEvent
}

func NewUpdatingTenant(tenant contracts.Tenant) *UpdatingTenant {
	return eventcontracts.NewTenantEvent("UpdatingTenant", tenant).(*UpdatingTenant)
}
