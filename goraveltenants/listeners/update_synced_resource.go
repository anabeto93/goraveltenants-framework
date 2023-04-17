package listeners

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database"
	"github.com/anabeto93/goraveltenants/database/models"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/anabeto93/goraveltenants/exceptions"
	"github.com/goravel/framework/facades"
	"reflect"
)

type UpdateSyncedResource struct {
	Database    *database.DatabaseManager
	shouldQueue bool
}

func NewUpdateSyncedResource(databaseManager *database.DatabaseManager) *UpdateSyncedResource {
	return &UpdateSyncedResource{
		Database:    databaseManager,
		shouldQueue: false,
	}
}

func (u *UpdateSyncedResource) Handle(args ...interface{}) error {
	syncedResourceEvent := args[0].(events.SyncedResourceSaved)
	syncedModel := syncedResourceEvent.GetModel()
	tenant := syncedResourceEvent.GetTenant()

	syncedAttributes := syncedModel.GetSyncedAttributeNames()

	var tenants []contracts.Tenant
	if tenant != nil {
		// Update resource in central database and get tenants
		tenants = u.updateResourceInCentralDatabaseAndGetTenants(syncedModel, tenant)
	} else {
		temp, err := u.getTenantsForCentralModel(syncedModel)
		if err != nil {
			return err
		}

		tenants = temp
	}

	// Update resource in tenant databases
	u.updateResourceInTenantDatabases(tenants, syncedModel, syncedAttributes)

	return nil
}

func (u *UpdateSyncedResource) Name() string {
	return "UpdateSyncedResource"
}

func (u *UpdateSyncedResource) convertToInterfaceSlice(tenants []models.Tenant) []contracts.Tenant {
	instanceTenants := make([]contracts.Tenant, len(tenants))

	for i, currentTenant := range tenants {
		instanceTenants[i] = &currentTenant
	}

	return instanceTenants
}

func (u *UpdateSyncedResource) getTenantsForCentralModel(syncedModel interface{}) ([]contracts.Tenant, error) {
	centralModel, ok := syncedModel.(*contracts.SyncMaster)
	if !ok {
		return nil, exceptions.NewModelNotSyncMasterException(reflect.TypeOf(syncedModel).String())
	}
	var tenants []models.Tenant
	if err := facades.Orm.Query().Model(centralModel).With("tenants").Load(&tenants, "tenants"); err != nil {
		return nil, err
	}

	return u.convertToInterfaceSlice(tenants), nil
}

func (u *UpdateSyncedResource) updateResourceInCentralDatabaseAndGetTenants(syncedModel contracts.Syncable, tenant contracts.Tenant) []contracts.Tenant {
	// Implement logic to update the resource in the central database and return tenants
	
}

func (u *UpdateSyncedResource) updateResourceInTenantDatabases(tenants []contracts.Tenant, syncedModel contracts.Syncable, syncedAttributes []string) {
	// Implement logic to update the resource in tenant databases
}
