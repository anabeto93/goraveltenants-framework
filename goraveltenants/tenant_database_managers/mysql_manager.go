package tenant_database_managers

import (
	"database/sql"
	"fmt"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/exceptions"
	"github.com/goravel/framework/facades"
)

var _ contracts.TenantDatabaseManager = &MySQLDatabaseManager{}

type MySQLDatabaseManager struct {
	connection   string
	dbConnection *sql.DB
}

func (m *MySQLDatabaseManager) SetConnection(connection string) error {
	m.connection = connection
	return nil
}

func (m *MySQLDatabaseManager) database() (*sql.DB, error) {
	if m.connection == "" {
		return nil, exceptions.NewNoConnectionSetException("MySQLDatabaseManager")
	}

	if m.dbConnection != nil {
		return m.dbConnection, nil
	}

	connection, err := facades.Orm.Connection(m.connection).DB()
	if err != nil {
		return nil, err
	}

	m.dbConnection = connection

	return connection, nil
}

func (m *MySQLDatabaseManager) getDBInfo(db *sql.DB) (string, string, string, error) {
	var dbName, charset, collation string

	// Execute query to get database info
	err := db.QueryRow("SELECT DATABASE() AS dbname, DEFAULT_CHARACTER_SET_NAME AS charset, DEFAULT_COLLATION_NAME AS collation FROM information_schema.SCHEMATA WHERE SCHEMA_NAME=?", dbName).Scan(&dbName, &charset, &collation)
	if err != nil {
		return "", "", "", err
	}

	return dbName, charset, collation, nil
}

func (m *MySQLDatabaseManager) CreateDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	database := tenant.Database().GetName()
	_, charset, collation, _ := m.getDBInfo(db)
	sqlStmt := fmt.Sprintf(`CREATE DATABASE %s CHARACTER SET %s COLLATE %s`, database, charset, collation)
	_, err = db.Exec(sqlStmt)
	return err == nil
}

func (m *MySQLDatabaseManager) DeleteDatabase(tenant contracts.TenantWithDatabase) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	sqlStmt := fmt.Sprintf(`DROP DATABASE %s`, tenant.Database().GetName())
	_, err = db.Exec(sqlStmt)
	return err == nil
}

func (m *MySQLDatabaseManager) DatabaseExists(name string) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	sqlStmt := fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", name)
	var dbName string
	err = db.QueryRow(sqlStmt).Scan(&dbName)
	if err != nil {
		return false
	}

	// Get the database's default charset and collation
	sqlStmt = fmt.Sprintf("SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", name)
	var charset, collation string
	err = db.QueryRow(sqlStmt).Scan(&charset, &collation)
	if err != nil {
		return false
	}

	return true
}

func (m *MySQLDatabaseManager) MakeConnectionConfig(baseConfig map[string]interface{}, databaseName string) map[string]interface{} {
	baseConfig["database"] = databaseName
	return baseConfig
}
