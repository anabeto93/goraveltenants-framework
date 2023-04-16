package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DatabaseDeleted{}

type DatabaseDeleted struct {
	eventcontracts.BaseTenantEvent
}

func NewDatabaseDeleted(tenant contracts.Tenant) *DatabaseDeleted {
	return eventcontracts.NewTenantEvent("DatabaseDeleted", tenant).(*DatabaseDeleted)
}
