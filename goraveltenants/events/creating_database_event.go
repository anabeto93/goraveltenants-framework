package events

import (
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
	"github.com/anabeto93/goraveltenants/contracts"
)

var _ eventcontracts.TenantEvent = &CreatingDatabase{}

type CreatingDatabase struct {
	eventcontracts.BaseTenantEvent
}

func NewCreatingDatabase(tenant contracts.Tenant) *CreatingDatabase {
	return eventcontracts.NewTenantEvent("CreatingDatabase", tenant).(*CreatingDatabase)
}
