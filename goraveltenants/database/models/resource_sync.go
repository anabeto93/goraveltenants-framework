package models

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

var _ contracts.ResourceSyncer = &BaseResourceSyncing{}

type BaseResourceSyncing struct{}

func (brs *BaseResourceSyncing) GetGlobalIdentifierKeyName() string {
	return "id"
}

func (brs *BaseResourceSyncing) TriggerSyncEvent() {
	_ = facades.Event.Job(events.NewSyncedResourceSaved(currentEvent.GetTenant()), []event.Arg{}).Dispatch()
}