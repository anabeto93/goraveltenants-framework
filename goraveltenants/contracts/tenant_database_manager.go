package contracts

// TenantDatabaseManager is an interface that defines the methods required to create and delete databases for tenants.
type TenantDatabaseManager interface {
	// CreateDatabase creates a database for the given tenant.
	CreateDatabase(tenant TenantWithDatabase) bool

	// DeleteDatabase deletes a database for the given tenant.
	DeleteDatabase(tenant TenantWithDatabase) bool

	// DatabaseExists checks whether a database with the given name exists.
	DatabaseExists(name string) bool

	// MakeConnectionConfig creates a DB connection config array for the given database name.
	MakeConnectionConfig(baseConfig map[string]interface{}, databaseName string) map[string]interface{}

	// SetConnection sets the DB connection that should be used by the tenant database manager.
	SetConnection(connection string) error
}
