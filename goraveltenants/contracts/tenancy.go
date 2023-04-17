package contracts

type Tenancy interface {
	Tenant
	Initialize(tenant interface{}) error
	End() error
	GetBootstrappers() []TenancyBootstrapper
	RunForMultiple(tenants []interface{}, callback func(tenant Tenant) error) error
	GetCurrentTenant(key ...interface{}) Tenant
}
