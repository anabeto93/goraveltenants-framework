package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &SavingTenant{}

type SavingTenant struct {
	eventcontracts.BaseTenantEvent
}

func NewSavingTenant(tenant contracts.Tenant) *SavingTenant {
	return eventcontracts.NewTenantEvent("SavingTenant", tenant).(*SavingTenant)
}
