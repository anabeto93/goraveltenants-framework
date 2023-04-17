package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/contracts/event"
)

type SyncedResourceSaved struct {
	model  contracts.Syncable
	tenant contracts.TenantWithDatabase
}

func (sr *SyncedResourceSaved) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

func (sr *SyncedResourceSaved) Name() string {
	return "SyncedResourceSaved"
}

func NewSyncedResourceSaved(model contracts.Syncable, tenant contracts.TenantWithDatabase) *SyncedResourceSaved {
	return &SyncedResourceSaved{
		model:  model,
		tenant: tenant,
	}
}

func (sr *SyncedResourceSaved) GetModel() contracts.Syncable {
	return sr.model
}

func (sr *SyncedResourceSaved) GetTenant() contracts.TenantWithDatabase {
	return sr.tenant
}
