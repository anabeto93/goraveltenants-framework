package tenant_database_managers

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
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
	hostname := databaseConfig.Connection()["host"]
	password := databaseConfig.GetPassword()

	m.connection.Exec(fmt.Sprintf("CREATE USER `%s`@`%%` IDENTIFIED BY '%s'", username, password))

	grantList := strings.Join(grants, ", ")

	var grantQuery string
	if isVersion8(m.connection) {
		grantQuery = fmt.Sprintf("GRANT %s ON `%s`.* TO `%s`@`%%`", grantList, database, username)
	} else {
		grantQuery = fmt.Sprintf("GRANT %s ON `%s`.* TO `%s`@`%%` IDENTIFIED BY '%s'", grantList, database, username, password)
	}

	return m.connection.Exec(grantQuery).Error == nil
}

func isVersion8(db *gorm.DB) bool {
	var version string
	db.Raw("SELECT VERSION()").Row().Scan(&version)

	return strings.HasPrefix(version, "8.")
}

func (m *PermissionControlledMySQLDatabaseManager) DeleteUser(databaseConfig contracts.DatabaseConfig) bool {
	username := databaseConfig.GetUsername()
	return m.connection.Exec(fmt.Sprintf("DROP USER IF EXISTS `%s`@`%%`", username)).Error == nil
}

func (m *PermissionControlledMySQLDatabaseManager) UserExists(username string) bool {
	var count int
	m.connection.Raw(fmt.Sprintf("SELECT COUNT(*) FROM mysql.user WHERE user = '%s'", username)).Row().Scan(&count)
	return count > 0
}
