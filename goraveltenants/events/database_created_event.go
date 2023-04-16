package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DatabaseCreated{}

type DatabaseCreated struct {
	eventcontracts.BaseTenantEvent
}

func NewDatabaseCreated(tenant contracts.Tenant) *DatabaseCreated {
	return eventcontracts.NewTenantEvent("DatabaseCreated", tenant).(*DatabaseCreated)
}
