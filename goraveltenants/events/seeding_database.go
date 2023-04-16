package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &SeedingDatabase{}

type SeedingDatabase struct {
	eventcontracts.BaseTenantEvent
}

func NewSeedingDatabase(tenant contracts.Tenant) *SeedingDatabase {
	return eventcontracts.NewTenantEvent("SeedingDatabase", tenant).(*SeedingDatabase)
}
