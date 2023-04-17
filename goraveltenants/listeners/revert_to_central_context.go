package listeners

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

var _ contracts.Listener = &RevertToCentralContext{}

type RevertToCentralContext struct {
	event events.TenancyEnded
}

func (bt *RevertToCentralContext) Handle(args ...interface{}) error {
	currentEvent := args[0].(events.TenancyEnded)
	bt.event = currentEvent

	_ = facades.Event.Job(events.NewRevertingToCentralContextEvent(currentEvent.GetTenant()), []event.Arg{}).Dispatch()

	for _, bootstrapper := range currentEvent.GetTenant().GetBootstrappers() {
		bootstrapper.Revert()
	}

	_ = facades.Event.Job(events.NewRevertedToCentralContextEvent(currentEvent.GetTenant()), []event.Arg{}).Dispatch()
	return nil
}

func (bt *RevertToCentralContext) Name() string {
	return "RevertToCentralContext"
}
