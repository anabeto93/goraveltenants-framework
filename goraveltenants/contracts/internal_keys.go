package contracts

type HasInternalKeys interface {
	// GetInternal returns the value of an internal key.
	GetInternal(key string) interface{}

	// SetInternal sets the value of an internal key.
	SetInternal(key string, value interface{})

	// GetAttributes returns the map[string]interface{} of all the models attributes
	GetAttributes() map[string]interface{}

	// InternalPrefix used to Prefix to all internal properties
	InternalPrefix() string
}
