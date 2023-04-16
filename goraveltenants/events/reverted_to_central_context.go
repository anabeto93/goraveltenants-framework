package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &RevertedToCentralContext{}

type RevertedToCentralContext struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *RevertedToCentralContext) Name() string {
	return "RevertedToCentralContext"
}

func NewRevertedToCentralContextEvent(tenant contracts.Tenant) *RevertedToCentralContext {
	return eventcontracts.NewTenancyEvent("RevertedToCentralContext", tenant).(*RevertedToCentralContext)
}
