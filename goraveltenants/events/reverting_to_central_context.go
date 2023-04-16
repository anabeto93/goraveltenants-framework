package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.TenancyEvent = &RevertingToCentralContext{}

type RevertingToCentralContext struct {
	eventcontracts.BaseTenancyEvent
}

func (bt *RevertingToCentralContext) Name() string {
	return "RevertingToCentralContext"
}

func NewRevertingToCentralContextEvent(tenant contracts.Tenant) *RevertingToCentralContext {
	return eventcontracts.NewTenancyEvent("RevertingToCentralContext", tenant).(*RevertingToCentralContext)
}
