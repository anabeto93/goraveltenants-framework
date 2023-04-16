package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DatabaseMigrated{}

type DatabaseMigrated struct {
	eventcontracts.BaseTenantEvent
}

func NewDatabaseMigrated(tenant contracts.Tenant) *DatabaseMigrated {
	return eventcontracts.NewTenantEvent("DatabaseMigrated", tenant).(*DatabaseMigrated)
}
