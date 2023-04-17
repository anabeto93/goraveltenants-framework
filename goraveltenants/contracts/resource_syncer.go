package contracts

type ResourceSyncer interface {
	GetGlobalIdentifierKeyName() string
	TriggerSyncEvent()
}
