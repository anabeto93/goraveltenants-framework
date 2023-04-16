package events

import "github.com/anabeto93/goraveltenants/contracts"

type SyncedResourceSaved struct {
	model  contracts.Syncable
	tenant contracts.TenantWithDatabase
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
