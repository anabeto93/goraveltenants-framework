package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenantEvent = &DeletingDatabase{}

type DeletingDatabase struct {
	eventcontracts.BaseTenantEvent
}

func NewDeletingDatabase(tenant contracts.Tenant) *DeletingDatabase {
	return eventcontracts.NewTenantEvent("DeletingDatabase", tenant).(*DeletingDatabase)
}
