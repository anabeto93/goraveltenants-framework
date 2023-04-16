package tenant_database_managers

import (
	"database/sql"
	"fmt"
	"github.com/goravel/framework/facades"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
)

var _ contracts.TenantDatabaseManager = &PostgresDatabaseManager{}

type PostgresDatabaseManager struct {
	connection   string
	dbConnection *sql.DB
}

func (p *PostgresDatabaseManager) SetConnection(connection string) error {
	p.connection = connection
	return nil
}

func (p *PostgresDatabaseManager) database() (*sql.DB, error) {
	if p.connection == "" {
		return nil, exceptions.NewNoConnectionSetException("PostgresDatabaseManager")
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

func (p *PostgresDatabaseManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sqlStmt := fmt.Sprintf(`CREATE DATABASE "%s"  WITH TEMPLATE=template0`, tenant.Database().GetName())
	_, err = db.Exec(sqlStmt)
	return err == nil
}

func (p *PostgresDatabaseManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := p.database()
	if err != nil {
		return false
	}
	sqlStmt := fmt.Sprintf(`DROP DATABASE "%s"`, tenant.Database().GetName())
	_, err = db.Exec(sqlStmt)
	return err == nil
}

func (p *PostgresDatabaseManager) DatabaseExists(name string) bool {
	db, err := p.database()
	if err != nil {
		return false
	}

	sqlStmt := fmt.Sprintf("SELECT datname FROM pg_database WHERE datname = '%s'", name)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return false
	}

	return rows.Next()
}

func (p *PostgresDatabaseManager) MakeConnectionConfig(baseConfig map[string]interface{}, databaseName string) map[string]interface{} {
	baseConfig["database"] = databaseName
	return baseConfig
}
