package goraveltenants

import (
	"database/sql"
	"errors"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/goravel/framework/contracts/event"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database/models"
	"github.com/goravel/framework/facades"
)

var _ contracts.Tenancy = &Tenancy{}

type Tenancy struct {
	models.Tenant
	initialized bool
}

func NewTenancy() *Tenancy {
	return &Tenancy{
		Tenant:      models.Tenant{},
		initialized: false,
	}
}

func (t *Tenancy) getDatabaseInstance() (*sql.DB, error) {
	conn := facades.Config.GetString("database.default")

	connection, err := facades.Orm.Connection(conn).DB()
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (t *Tenancy) Initialize(tenant interface{}) error {
	if tenant == nil {
		return errors.New("tenant cannot be nil")
	}

	var tenantInstance models.Tenant
	switch v := tenant.(type) {
	case models.Tenant:
		tenantInstance = v
	case string, int:
		tenantFound, err := t.findTenantByKey(tenant)
		if err != nil {
			return err
		}
		tenantInstance = tenantFound
	default:
		return errors.New("invalid tenant type")
	}

	if t.initialized && t.Tenant.GetTenantKey() == tenantInstance.GetTenantKey() {
		return nil
	}

	if t.initialized {
		if err := t.End(); err != nil {
			return err
		}
	}

	t.Tenant = tenantInstance
	facades.Event.Job(events.NewInitializingTenancyEvent(t), []event.Arg{})

	t.initialized = true

	facades.Event.Job(events.NewTenancyInitializedEvent(t), []event.Arg{})

	return nil
}

func (t *Tenancy) findTenantByKey(tenant interface{}) (models.Tenant, error) {
	tenantModel := models.Tenant{}
	// Replace with actual logic to find the tenant from the database
	sqlDB, err := t.getDatabaseInstance()
	if err != nil {
		return models.Tenant{}, err
	}

	query := "SELECT * FROM " + tenantModel.TableName() + " WHERE " + tenantModel.GetTenantKeyName() + " = ?"
	if err := sqlDB.QueryRow(query, tenant).Scan(&tenantModel.ID, &tenantModel.CreatedAt, &tenantModel.UpdatedAt, &tenantModel.DeletedAt, &tenantModel.Data); err != nil {
		return models.Tenant{}, err
	}
	return tenantModel, nil
}

func (t *Tenancy) End() error {
	// Emit events here if necessary

	if !t.initialized {
		return nil
	}

	// Emit events here if necessary
	facades.Event.Job(events.NewEndingTenancyEvent(t), []event.Arg{})
	t.initialized = false
	t.Tenant = models.Tenant{}
	facades.Event.Job(events.NewTenancyEndedEvent(t), []event.Arg{})

	return nil
}

func (t *Tenancy) GetBootstrappers() []contracts.TenancyBootstrapper {
	return t.GetConfig("tenancy.bootstrappers").([]contracts.TenancyBootstrapper)
}

func (t *Tenancy) RunForMultiple(tenants []interface{}, callback func(tenant contracts.Tenant) error) error {
	originalTenant := t.Tenant

	for _, tenant := range tenants {
		var tenantInstance models.Tenant
		switch v := tenant.(type) {
		case models.Tenant:
			tenantInstance = v
		case string, int:
			tenantFound, err := t.findTenantByKey(tenant)
			if err != nil {
				return err
			}
			tenantInstance = tenantFound
		default:
			return errors.New("invalid tenant type")
		}

		err := t.Initialize(tenantInstance)
		if err != nil {
			return err
		}

		err = callback(&tenantInstance)
		if err != nil {
			return err
		}
	}

	if originalTenant != (models.Tenant{}) {
		if err := t.Initialize(originalTenant); err != nil {
			return err
		}
	} else {
		if err := t.End(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tenancy) GetConfig(key string) interface{} {
	return facades.Config.Get(key)
}

func (t *Tenancy) GetCurrentTenant(key ...interface{}) contracts.Tenant {
	if t.initialized == true {
		return &t.Tenant
	}

	// first check that the database.connections.tenant exists
	tempConn := facades.Config.Get("database.connections.tenant")
	if tempConn == nil {
		return nil // it just means there's currently no tenant available
	}
	// In our DatabaseManager.CreateTenantConnection(), this is supposed to be a map[string]interface{}
	if _, ok := tempConn.(map[string]interface{}); !ok {
		return nil
	}
	// by this point it means we might just have a tenant available and currently setup and connected to it
	sqlDB, err := t.getDatabaseInstance()
	if err != nil {
		return nil
	}
}
