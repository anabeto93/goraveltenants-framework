package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &MigratingDatabase{}

type MigratingDatabase struct {
	eventcontracts.BaseTenantEvent
}

func NewMigratingDatabase(tenant contracts.Tenant) *MigratingDatabase {
	return eventcontracts.NewTenantEvent("MigratingDatabase", tenant).(*MigratingDatabase)
}
