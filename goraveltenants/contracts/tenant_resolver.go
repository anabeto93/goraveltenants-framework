package contracts

type TenantResolver interface {
	Resolve(args ...interface{}) (Tenant, TenantCouldNotBeIdentifiedException)
}
