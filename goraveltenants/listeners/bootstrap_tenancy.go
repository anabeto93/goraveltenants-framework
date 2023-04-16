package listeners

import (
	"github.com/anabeto93/goraveltenants/events"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type BootstrapTenancy struct {
	event events.TenancyInitialized
}

func (bt *BootstrapTenancy) Handle(args ...interface{}) error {
	currentEvent := args[0].(events.TenancyInitialized)
	bt.event = currentEvent

	_ = facades.Event.Job(events.NewBootstrappingTenancyEvent(currentEvent.GetTenant()), []event.Arg{}).Dispatch()

	for _, bootstrapper := range currentEvent.GetTenant().GetBootstrappers() {
		bootstrapper.Bootstrap(currentEvent.GetTenant())
	}

	_ = facades.Event.Job(events.NewTenancyBootstrappedEvent(currentEvent.GetTenant()), []event.Arg{}).Dispatch()
	return nil
}
