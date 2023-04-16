package tenant_database_managers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/anabeto93/goraveltenants/contracts"
)

var _ contracts.ManagesDatabaseUsers = &PermissionControlledMySQLDatabaseManager{}
var _ contracts.TenantDatabaseManager = &PermissionControlledMySQLDatabaseManager{}

type PermissionControlledMySQLDatabaseManager struct {
	MySQLDatabaseManager
}

var grants = []string{
	"ALTER", "ALTER ROUTINE", "CREATE", "CREATE ROUTINE", "CREATE TEMPORARY TABLES", "CREATE VIEW",
	"DELETE", "DROP", "EVENT", "EXECUTE", "INDEX", "INSERT", "LOCK TABLES", "REFERENCES", "SELECT",
	"SHOW VIEW", "TRIGGER", "UPDATE",
}

func (m *PermissionControlledMySQLDatabaseManager) CreateUser(databaseConfig contracts.DatabaseConfig) bool {
	database := databaseConfig.GetName()
	username := databaseConfig.GetUsername()
	password := databaseConfig.GetPassword()

	//m.connection.Exec(fmt.Sprintf("CREATE USER `%s`@`%%` IDENTIFIED BY '%s'", username, password))
	db, err := m.database()
	if err != nil {
		return false
	}
	if _, err = db.Exec(fmt.Sprintf("CREATE USER `%s`@`%%` IDENTIFIED BY '%s'", username, password)); err != nil {
		return false
	}

	grantList := strings.Join(grants, ", ")

	var grantQuery string
	if isVersion8(db) {
		grantQuery = fmt.Sprintf("GRANT %s ON `%s`.* TO `%s`@`%%`", grantList, database, username)
	} else {
		grantQuery = fmt.Sprintf("GRANT %s ON `%s`.* TO `%s`@`%%` IDENTIFIED BY '%s'", grantList, database, username, password)
	}

	_, err = db.Exec(grantQuery)
	return err == nil
}

func isVersion8(db *sql.DB) bool {
	var version string
	err := db.QueryRow("SELECT VERSION();").Scan(&version)
	if err != nil {
		// handle error
	}
	majorVersion := strings.Split(version, ".")[0]
	if majorVersion >= "8" {
		return true
	}
	return false
}

func (m *PermissionControlledMySQLDatabaseManager) DeleteUser(databaseConfig contracts.DatabaseConfig) bool {
	username := databaseConfig.GetUsername()
	db, err := m.database()
	if err != nil {
		return false
	}

	_, err = db.Exec(fmt.Sprintf("DROP USER IF EXISTS `%s`@`%%`", username))
	return err == nil
}

func (m *PermissionControlledMySQLDatabaseManager) UserExists(username string) bool {
	db, err := m.database()
	if err != nil {
		return false
	}
	var count int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM mysql.user WHERE user = '%s'", username)).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}
