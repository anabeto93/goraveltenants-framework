package contracts

type Syncable interface {
	GetGlobalIdentifierKeyName() string
	GetGlobalIdentifierKey() string
	GetCentralModelName() string
	GetSyncedAttributeNames() []string
	TriggerSyncEvent()
}
