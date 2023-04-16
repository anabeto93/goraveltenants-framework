package tenant_database_managers

import (
	"database/sql"
	"fmt"
	"github.com/goravel/framework/facades"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
)

var _ contracts.TenantDatabaseManager = &PostgreSQLSchemaManager{}

type PostgreSQLSchemaManager struct {
	connection   string
	dbConnection *sql.DB
}

func (p *PostgreSQLSchemaManager) SetConnection(connection string) error {
	p.connection = connection
	return nil
}

func (p *PostgreSQLSchemaManager) database() (*sql.DB, error) {
	if p.connection == "" {
		return nil, exceptions.NewNoConnectionSetException("PostgreSQLSchemaManager")
	}

	if p.dbConnection != nil {
		return p.dbConnection, nil
	}

	connection, err := facades.Orm.Connection(p.connection).DB()
	if err != nil {
		return nil, err
	}

	p.dbConnection = connection

	return connection, nil
}

func (p *PostgreSQLSchemaManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sqlStmt := fmt.Sprintf(`CREATE SCHEMA "%s"`, tenant.Database().GetName())
	_, err = db.Exec(sqlStmt)
	return err == nil
}

func (p *PostgreSQLSchemaManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sqlStmt := fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, tenant.Database().GetName())
	_, err = db.Exec(sqlStmt)
	return err == nil
}

func (p *PostgreSQLSchemaManager) DatabaseExists(name string) bool {
	db, err := p.database()
	if err != nil {
		return false
	}

	sqlStmt := fmt.Sprintf("SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s'", name)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return false
	}

	return rows.Next()
}

func (p *PostgreSQLSchemaManager) MakeConnectionConfig(baseConfig map[string]interface{}, databaseName string) map[string]interface{} {
	baseConfig["search_path"] = databaseName
	return baseConfig
}
