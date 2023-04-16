package goraveltenants

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
	"github.com/anabeto93/goraveltenants/tenant_database_managers"
	"github.com/goravel/framework/facades"
	"strings"
)

var _ contracts.DatabaseConfig = &DatabaseConfig{}

type DatabaseConfig struct {
	tenant contracts.Tenant
}

func NewDatabaseConfig(tenant contracts.Tenant) *DatabaseConfig {
	return &DatabaseConfig{tenant: tenant}
}

func (dc *DatabaseConfig) GetName() string {
	return dc.tenant.GetInternal("db_name")
}

func (dc *DatabaseConfig) GetUsername() string {
	return dc.tenant.GetInternal("db_username")
}

func (dc *DatabaseConfig) GetPassword() string {
	return dc.tenant.GetInternal("db_password")
}

func (dc *DatabaseConfig) MakeCredentials() {
	dc.tenant.SetInternal("db_name", dc.GetName())

	if manager, ok := dc.Manager().(contracts.ManagesDatabaseUsers); ok {
		dc.tenant.SetInternal("db_username", manager.CreateUser(dc))
		dc.tenant.SetInternal("db_password", manager.CreateUser(dc))
	}
}

func (dc *DatabaseConfig) GetTemplateConnectionName() string {
	return dc.tenant.GetInternal("db_connection")
}

func (dc *DatabaseConfig) Connection() map[string]interface{} {
	template := dc.GetTemplateConnectionName()
	templateConnection := facades.Config.GetString("database.connections." + template)

	return dc.Manager().MakeConnectionConfig(templateConnection, dc.GetName())
}

func (dc *DatabaseConfig) TenantConfig() map[string]interface{} {
	keys := []string{}
	for key := range dc.tenant.internalData {
		if strings.HasPrefix(key, "db_") {
			keys = append(keys, key)
		}
	}

	result := make(map[string]interface{})
	for _, key := range keys {
		result[strings.TrimPrefix(key, "db_")] = dc.tenant.GetInternal(key)
	}

	return result
}

func (dc *DatabaseConfig) Manager() contracts.TenantDatabaseManager {
	driver := facades.Config.GetString("database.connections." + dc.GetTemplateConnectionName() + ".driver")
	var databaseManagers map[string]string
	databaseManagers = facades.Config.Get("tenancy.database.managers").(map[string]string)

	managerName, ok := databaseManagers[driver]
	if !ok {
		err := exceptions.NewDatabaseManagerNotRegisteredException(driver)
		panic(err.Error())
	}

	var manager contracts.TenantDatabaseManager

	switch managerName {
	case "mysql":
		manager = &tenant_database_managers.MySQLDatabaseManager{}
	case "pgsql":
		manager = &tenant_database_managers.PostgreSQLSchemaManager{}
	case "mysql_permissions":
		manager = &tenant_database_managers.PermissionControlledMySQLDatabaseManager{}
	case "sqlite":
		manager = &tenant_database_managers.SQLiteDatabaseManager{}
	}

	manager.SetConnection(dc.GetTemplateConnectionName())

	return manager
}
