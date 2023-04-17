package database

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/anabeto93/goraveltenants/contracts"
)

type DatabaseManager struct {
	connections map[string]*sql.DB
	mu          sync.Mutex
}

func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		connections: make(map[string]*sql.DB),
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
			conn.Close()
			delete(dm.connections, "tenant")
		}
		dm.mu.Unlock()
	}
		
	facades.Config.Add("database.connections.tenant", nil)
}

func (dm *DatabaseManager) EnsureTenantCanBeCreated(tenant contracts.TenantWithDatabase) error {
	manager := tenant.Database().Manager()

	if manager.DatabaseExists(database := tenant.Database().GetName()) {
		return &TenantDatabaseAlreadyExistsException{DatabaseName: database}
	}

	if userMgr, ok := manager.(contracts.ManagesDatabaseUsers); ok {
		if userMgr.UserExists(username := tenant.Database().GetUsername()) {
			return &TenantDatabaseUserAlreadyExistsException{Username: username}
		}
	}

	return nil
}
