package models

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/anabeto93/goraveltenants/facades"
	"github.com/goravel/framework/contracts/event"
	frameworkfacades "github.com/goravel/framework/facades"
)

func TriggerSyncEvent(model contracts.Syncable) {
	currentTenant := facades.Tenancy.GetCurrentTenant()
	tenantWithDb, ok := currentTenant.(contracts.TenantWithDatabase)
	if !ok {
		tenantWithDb = NewBaseTenantWithDatabase(currentTenant.(*Tenant))
	}
	_ = frameworkfacades.Event.Job(events.NewSyncedResourceSaved(model, tenantWithDb), []event.Arg{}).Dispatch()
}
