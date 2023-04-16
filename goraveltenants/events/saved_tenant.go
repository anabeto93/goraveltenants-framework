package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &SavedTenant{}

type SavedTenant struct {
	eventcontracts.BaseTenantEvent
}

func NewSavedTenant(tenant contracts.Tenant) *SavedTenant {
	return eventcontracts.NewTenantEvent("SavedTenant", tenant).(*SavedTenant)
}
