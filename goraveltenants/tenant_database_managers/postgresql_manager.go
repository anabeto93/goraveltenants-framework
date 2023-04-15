package tenant_database_managers

import (
	"fmt"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
	"gorm.io/gorm"
)

var _ contracts.TenantDatabaseManager = &PostgreSQLSchemaManager{}

type PostgreSQLSchemaManager struct {
	connection *gorm.DB
}

func (p *PostgreSQLSchemaManager) SetConnection(connection *gorm.DB) error {
	p.connection = connection
	return nil
}

func (p *PostgreSQLSchemaManager) database() (*gorm.DB, error) {
	if p.connection == nil {
		return nil, exceptions.NewNoConnectionSetException("PostgreSQLSchemaManager")
	}
	return p.connection, nil
}

func (p *PostgreSQLSchemaManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf(`CREATE SCHEMA "%s"`, tenant.Database().GetName())
	return db.Exec(sql).Error == nil
}

func (p *PostgreSQLSchemaManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, tenant.Database().GetName())
	return db.Exec(sql).Error == nil
}

func (p *PostgreSQLSchemaManager) DatabaseExists(name string) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sql := fmt.Sprintf("SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s'", name)
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (p *PostgreSQLSchemaManager) MakeConnectionConfig(baseConfig map[string]string, databaseName string) map[string]string {
	baseConfig["search_path"] = databaseName
	return baseConfig
}
