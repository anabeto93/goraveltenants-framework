package tenant_database_managers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anabeto93/goraveltenants/contracts"
)

var _ contracts.TenantDatabaseManager = &SQLiteDatabaseManager{}

type SQLiteDatabaseManager struct {
	connection string
}

func (m *SQLiteDatabaseManager) SetConnection(connection string) error {
	m.connection = connection
	return nil
}

func databasePath(databaseName string) string {
	return filepath.Join("databases", fmt.Sprintf("%s.db", databaseName))
}

func (m *SQLiteDatabaseManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	file, err := os.Create(databasePath(tenant.Database().GetName()))
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

func (m *SQLiteDatabaseManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	err := os.Remove(databasePath(tenant.Database().GetName()))
	if err != nil {
		return false
	}
	return true
}

func (m *SQLiteDatabaseManager) DatabaseExists(name string) bool {
	_, err := os.Stat(databasePath(name))
	return !os.IsNotExist(err)
}

func (m *SQLiteDatabaseManager) MakeConnectionConfig(baseConfig map[string]interface{}, databaseName string) map[string]interface{} {
	baseConfig["database"] = databasePath(databaseName)
	return baseConfig
}
