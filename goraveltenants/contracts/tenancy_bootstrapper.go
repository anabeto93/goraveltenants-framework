package contracts

// TenancyBootstrapper is an interface that defines the methods required to make your application tenant-aware automatically.
type TenancyBootstrapper interface {
	// Bootstrap makes your application tenant-aware for the given tenant.
	Bootstrap(tenant Tenant)

	// Revert reverts the changes made by Bootstrap method.
	Revert()
}
