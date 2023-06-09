package database

import (
	"database/sql"
	"github.com/anabeto93/goraveltenants/exceptions"
	"sync"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/facades"
)

type DatabaseManager struct {
	connections map[string]interface{}
	mu          sync.Mutex
}

func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		connections: make(map[string]interface{}),
	}
}

func (dm *DatabaseManager) ConnectToTenant(tenant contracts.TenantWithDatabase) {
	dm.PurgeTenantConnection()
	dm.CreateTenantConnection(tenant)
	facades.Config.Add("database.default", "tenant")
}

func (dm *DatabaseManager) ReconnectToCentral() {
	dm.PurgeTenantConnection()
	facades.Config.Add("database.default", facades.Config.Get("tenancy.database.central_connection"))
}

func (dm *DatabaseManager) CreateTenantConnection(tenant contracts.TenantWithDatabase) {
	connection := tenant.Database().Connection()

	dm.mu.Lock()
	dm.connections["tenant"] = connection
	dm.mu.Unlock()

	facades.Config.Add("database.connections.tenant", connection)
}

func (dm *DatabaseManager) PurgeTenantConnection() {
	if exists := facades.Config.Get("database.connections.tenant"); exists != nil {
		dm.mu.Lock()
		if conn, ok := dm.connections["tenant"]; ok {
			if sqlConn, ok := conn.(*sql.DB); ok {
				_ = sqlConn.Close()
			}
			delete(dm.connections, "tenant")
		}
		dm.mu.Unlock()
	}

	facades.Config.Add("database.connections.tenant", nil)
}

func (dm *DatabaseManager) EnsureTenantCanBeCreated(tenant contracts.TenantWithDatabase) error {
	manager := tenant.Database().Manager()

	database := tenant.Database().GetName()
	if manager.DatabaseExists(database) {
		return exceptions.NewTenantDatabaseAlreadyExistsException(database)
	}

	if userMgr, ok := manager.(contracts.ManagesDatabaseUsers); ok {
		username := tenant.Database().GetUsername()
		if userMgr.UserExists(username) {
			return exceptions.NewTenantDatabaseUserAlreadyExistsException(username)
		}
	}

	return nil
}
