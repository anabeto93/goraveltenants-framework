package tenant_database_managers

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
)

var _ contracts.TenantDatabaseManager = &MySQLDatabaseManager{}

type MySQLDatabaseManager struct {
	connection *gorm.DB
}

func (m *MySQLDatabaseManager) SetConnection(connection *gorm.DB) error {
	m.connection = connection
	return nil
}

func (m *MySQLDatabaseManager) database() (*gorm.DB, error) {
	if m.connection == nil {
		return nil, exceptions.NewNoConnectionSetException("MySQLDatabaseManager")
	}
	return m.connection, nil
}

func (m *MySQLDatabaseManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	database := tenant.Database().GetName()
	charset := db.Dialect().GetConfig().Charset
	collation := db.Dialect().GetConfig().Collation
	sql := fmt.Sprintf(`CREATE DATABASE %s CHARACTER SET %s COLLATE %s`, database, charset, collation)
	return db.Exec(sql).Error == nil
}

func (m *MySQLDatabaseManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf(`DROP DATABASE %s`, tenant.Database().GetName())
	return db.Exec(sql).Error == nil
}

func (m *MySQLDatabaseManager) DatabaseExists(name string) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", name)
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (m *MySQLDatabaseManager) MakeConnectionConfig(baseConfig map[string]string, databaseName string) map[string]string {
	baseConfig["database"] = databaseName
	return baseConfig
}
