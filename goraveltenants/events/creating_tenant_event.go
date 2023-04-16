package events

import (
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
	"github.com/anabeto93/goraveltenants/contracts"
)

var _ eventcontracts.TenantEvent = &CreatingTenant{}

type CreatingTenant struct {
	eventcontracts.BaseTenantEvent
}

func NewCreatingTenant(tenant contracts.Tenant) *CreatingTenant {
	return eventcontracts.NewTenantEvent("CreatingTenant", tenant).(*CreatingTenant)
}
