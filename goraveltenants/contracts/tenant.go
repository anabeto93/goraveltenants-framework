package contracts

// Tenant is an interface that defines the methods required to implement a tenant model in Stancl Tenancy.
type Tenant interface {
	// GetTenantKeyName returns the name of the key used for identifying the tenant.
	GetTenantKeyName() string

	// GetTenantKey returns the value of the key used for identifying the tenant.
	GetTenantKey() interface{}

	// GetInternal returns the value of an internal key.
	GetInternal(key string) interface{}

	// SetInternal sets the value of an internal key.
	SetInternal(key string, value interface{})

	// Run runs a callback in this tenant's environment.
	Run(callback func())
}
