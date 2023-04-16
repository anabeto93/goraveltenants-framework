package tenant_database_managers

import (
	"fmt"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
	"gorm.io/gorm"
)

var _ contracts.TenantDatabaseManager = &PostgresDatabaseManager{}

type PostgresDatabaseManager struct {
	connection string
}

func (p *PostgresDatabaseManager) SetConnection(connection string) error {
	p.connection = connection
	return nil
}

func (p *PostgresDatabaseManager) database() (*gorm.DB, error) {
	if p.connection == nil {
		return nil, exceptions.NewNoConnectionSetException("PostgresDatabaseManager")
	}
	return p.connection, nil
}

func (p *PostgresDatabaseManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf(`CREATE DATABASE "%s"  WITH TEMPLATE=template0`, tenant.Database().GetName())
	return db.Exec(sql).Error == nil
}

func (p *PostgresDatabaseManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf(`DROP DATABASE "%s"`, tenant.Database().GetName())
	return db.Exec(sql).Error == nil
}

func (p *PostgresDatabaseManager) DatabaseExists(name string) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf("SELECT datname FROM pg_database WHERE datname = '%s'", name)
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (p *PostgresDatabaseManager) MakeConnectionConfig(baseConfig map[string]interface{}, databaseName string) map[string]interface{} {
	baseConfig["database"] = databaseName
	return baseConfig
}
