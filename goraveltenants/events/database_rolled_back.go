package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DatabaseRolledBack{}

type DatabaseRolledBack struct {
	eventcontracts.BaseTenantEvent
}

func NewDatabaseRolledBack(tenant contracts.Tenant) *DatabaseRolledBack {
	return eventcontracts.NewTenantEvent("DatabaseRolledBack", tenant).(*DatabaseRolledBack)
}
