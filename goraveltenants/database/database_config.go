package database

import (
	"github.com/anabeto93/goraveltenants"
	"strings"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
	"github.com/anabeto93/goraveltenants/tenant_database_managers"
	"github.com/goravel/framework/facades"
)

var _ contracts.DatabaseConfig = &DatabaseConfig{}

type DatabaseConfig struct {
	tenant                contracts.Tenant
	usernameGenerator     func(args ...interface{}) (string, error)
	passwordGenerator     func(args ...interface{}) (string, error)
	databaseNameGenerator func(args ...interface{}) (string, error)
}

func NewDatabaseConfig(tenant contracts.Tenant) *DatabaseConfig {
	return &DatabaseConfig{
		tenant: tenant,
		usernameGenerator: func(args ...interface{}) (string, error) {
			return goraveltenants.GenerateSecureRandomString(16)
		},
		passwordGenerator: func(args ...interface{}) (string, error) {
			return goraveltenants.GenerateSecureRandomString(32)
		},
		databaseNameGenerator: func(args ...interface{}) (string, error) {
			currentTenant := args[0].(contracts.Tenant)
			prefix := facades.Config.GetString("tenancy.database.prefix")
			key := currentTenant.GetTenantKey().(string)
			suffix := facades.Config.GetString("tenancy.database.suffix")
			return prefix + key + suffix, nil
		},
	}
}

func (dc *DatabaseConfig) GetName() string {
	name := dc.tenant.GetInternal("db_name").(string)

	if strings.TrimSpace(name) == "" {
		name, _ = dc.databaseNameGenerator(dc.tenant)
	}

	return name
}

func (dc *DatabaseConfig) GetUsername() string {
	name := dc.tenant.GetInternal("db_username").(string)

	if strings.TrimSpace(name) == "" {
		name, _ = dc.usernameGenerator(dc.tenant)
	}

	return name
}

func (dc *DatabaseConfig) GetPassword() string {
	password := dc.tenant.GetInternal("db_password").(string)

	if strings.TrimSpace(password) == "" {
		password, _ = dc.passwordGenerator(dc.tenant)
	}

	return password
}

func (dc *DatabaseConfig) MakeCredentials() {
	dc.tenant.SetInternal("db_name", dc.GetName())

	if _, ok := dc.Manager().(contracts.ManagesDatabaseUsers); ok {
		dc.tenant.SetInternal("db_username", dc.GetUsername())
		dc.tenant.SetInternal("db_password", dc.GetPassword())
	}
}

func (dc *DatabaseConfig) GetTemplateConnectionName() string {
	conn, ok := dc.tenant.GetInternal("db_connection").(string)
	if !ok {
		conn = ""
	}

	if strings.TrimSpace(conn) == "" {
		conn = facades.Config.GetString("tenancy.database.template_tenant_connection")
	}

	if strings.TrimSpace(conn) == "" {
		conn = facades.Config.GetString("tenancy.database.central_connection")
	}

	return conn
}

func (dc *DatabaseConfig) Connection() map[string]interface{} {
	template := dc.GetTemplateConnectionName()
	templateConnection := facades.Config.Get("database.connections." + template).(map[string]interface{})

	tenantConfig := dc.TenantConfig()

	config := MergeConfigMaps(tenantConfig, templateConnection)

	return dc.Manager().MakeConnectionConfig(config, dc.GetName())
}

func (dc *DatabaseConfig) TenantConfig() map[string]interface{} {
	var keys []string
	internalTenant, ok := dc.tenant.(contracts.HasInternalKeys)
	internalPrefix := "db_"
	if ok {
		internalPrefix = internalTenant.InternalPrefix() + "db_"
	}

	for key := range dc.tenant.GetAttributes() {
		if strings.HasPrefix(key, internalPrefix) {
			keys = append(keys, key)
		}
	}

	keys = removeKey(keys, internalPrefix+"name")
	keys = removeKey(keys, internalPrefix+"connection")

	result := make(map[string]interface{})
	for _, key := range keys {
		result[strings.TrimPrefix(key, internalPrefix)] = dc.tenant.GetInternal(key)
	}

	return result
}

/*func (dc *DatabaseConfig) TenantConfig() map[string]interface{} {
	var keys []string
	internalTenant, ok := dc.tenant.(contracts.HasInternalKeys)
	var result map[string]interface{}
	if !ok {
		for key, _ := range dc.tenant.GetAttributes() {
			if strings.HasPrefix(key, "db_") {
				keys = append(keys, key)
			}
		}

		keys = removeKey(keys, "db_name")
		keys = removeKey(keys, "db_connection")

		result = make(map[string]interface{})
		for _, key := range keys {
			result[strings.TrimPrefix(key, "db_")] = dc.tenant.GetInternal(key)
		}
	} else {
		for key, _ := range dc.tenant.GetAttributes() {
			if strings.HasPrefix(key, internalTenant.InternalPrefix()+"db_") {
				keys = append(keys, key)
			}
		}

		keys = removeKey(keys, internalTenant.InternalPrefix()+"db_name")
		keys = removeKey(keys, internalTenant.InternalPrefix()+"db_connection")

		result = make(map[string]interface{})
		for _, key := range keys {
			result[strings.TrimPrefix(key, internalTenant.InternalPrefix()+"db_")] = dc.tenant.GetInternal(key)
		}
	}

	return result
}*/

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

	_ = manager.SetConnection(dc.GetTemplateConnectionName())

	return manager
}

func (dc *DatabaseConfig) GeneratePasswordUsing(passwordGenerator func(args ...interface{}) (string, error)) {
	dc.passwordGenerator = passwordGenerator
}

func (dc *DatabaseConfig) GenerateDatabaseNameUsing(databaseNameGenerator func(args ...interface{}) (string, error)) {
	dc.databaseNameGenerator = databaseNameGenerator
}

func (dc *DatabaseConfig) GenerateUsernameUsing(usernameGenerator func(args ...interface{}) (string, error)) {
	dc.usernameGenerator = usernameGenerator
}

func MergeConfigMaps(first map[string]interface{}, second map[string]interface{}) map[string]interface{} {
	for key, val := range second {
		first[key] = val
	}

	return first
}

func removeKey(keys []string, keyToRemove string) []string {
	index := -1
	for i, key := range keys {
		if key == keyToRemove {
			index = i
			break
		}
	}

	if index != -1 {
		keys = append(keys[:index], keys[index+1:]...)
	}

	return keys
}
