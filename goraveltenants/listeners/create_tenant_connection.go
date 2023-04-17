package listeners

import (
	"github.com/anabeto93/goraveltenants/database"
	"github.com/anabeto93/goraveltenants/database/models"
	"github.com/anabeto93/goraveltenants/events/contracts"
)

type CreateTenantConnection struct {
	Database *database.DatabaseManager
}

func NewCreateTenantConnection(databaseManager *database.DatabaseManager) *CreateTenantConnection {
	return &CreateTenantConnection{
		Database: databaseManager,
	}
}

func (c *CreateTenantConnection) Handle(args ...interface{}) error {
	tenantEvent := args[0].(contracts.TenantEvent)
	tenant := tenantEvent.GetTenant()
	tenantWithDb := models.NewBaseTenantWithDatabase(tenant.(*models.Tenant))
	c.Database.CreateTenantConnection(tenantWithDb)
	return nil
}

func (c *CreateTenantConnection) Name() string {
	return "CreateTenantConnection"
}
