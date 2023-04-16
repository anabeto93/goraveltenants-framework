package events

import "github.com/anabeto93/goraveltenants/contracts"

type SyncedResourceChangedInForeignDatabase struct {
	model  contracts.Syncable
	tenant contracts.TenantWithDatabase
}

func (src *SyncedResourceChangedInForeignDatabase) Name() string {
	return "SyncedResourceChangedInForeignDatabase"
}

func NewSyncedResourceChangedInForeignDatabase(model contracts.Syncable, tenant contracts.TenantWithDatabase) *SyncedResourceChangedInForeignDatabase {
	return &SyncedResourceChangedInForeignDatabase{
		model:  model,
		tenant: tenant,
	}
}
