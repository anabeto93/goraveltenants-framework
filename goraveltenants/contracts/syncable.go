package contracts

type Syncable interface {
	GetGlobalIdentifierKeyName() string
	GetGlobalIdentifierKey() string
	GetCentralModelName() string
	GetCentralModelInstance() Syncable
	UpdateSyncedAttributesWithoutEvents(attributes map[string]interface{}) error
	GetSyncedAttributeNames() []string
	TriggerSyncEvent()
}
