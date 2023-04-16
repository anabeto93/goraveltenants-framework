package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DatabaseSeeded{}

type DatabaseSeeded struct {
	eventcontracts.BaseTenantEvent
}

func NewDatabaseSeeded(tenant contracts.Tenant) *DatabaseSeeded {
	return eventcontracts.NewTenantEvent("DatabaseSeeded", tenant).(*DatabaseSeeded)
}
