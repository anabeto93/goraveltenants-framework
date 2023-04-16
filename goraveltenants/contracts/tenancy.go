package contracts

type Tenancy interface {
	Tenant
	Initialize(tenant Tenant) error
	End() error
	GetBootstrappers() []TenancyBootstrapper
	RunForMultiple(tenants []Tenant, callback func(tenant Tenant) error) error
}
