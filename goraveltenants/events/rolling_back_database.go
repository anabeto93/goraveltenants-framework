package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &RollingBackDatabase{}

type RollingBackDatabase struct {
	eventcontracts.BaseTenantEvent
}

func NewRollingBackDatabase(tenant contracts.Tenant) *RollingBackDatabase {
	return eventcontracts.NewTenantEvent("RollingBackDatabase", tenant).(*RollingBackDatabase)
}
