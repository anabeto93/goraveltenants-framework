package tenant_database_managers

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
	"github.com/anabeto93/goraveltenants/contracts"
)

type SQLiteDatabaseManager struct {
	connection *gorm.DB
}

func (m *SQLiteDatabaseManager) SetConnection(connection *gorm.DB) error {
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
	file.Close()
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

func (m *SQLiteDatabaseManager) MakeConnectionConfig(baseConfig map[string]string, databaseName string) map[string]string {
	baseConfig["database"] = databasePath(databaseName)
	return baseConfig
}
