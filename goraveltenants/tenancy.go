package goraveltenants

import (
	"database/sql"
	"errors"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database/models"
	"github.com/goravel/framework/facades"
)

var _ contracts.Tenancy = &Tenancy{}

type Tenancy struct {
	models.Tenant
	initialized bool
}

func (t *Tenancy) getDatabaseInstance() (*sql.DB, error) {
	conn := facades.Config.GetString("database.default")

	connection, err := facades.Orm.Connection(conn).DB()
	if err != nil {
		return nil, err
	}

	return connection, nil
}

type TenantParam interface {
    ~int | ~string | contracts.Tenant
}

func (t *Tenancy) Initialize[T TenantParam](tenant T) error {
	if tenant == nil {
		return errors.New("tenant cannot be nil")
	}

	var tenantInstance models.Tenant
	switch v := tenant.(type) {
	case models.Tenant:
		tenantInstance = v
	case string, int:
		tenantModel := models.Tenant{}
		// Replace with actual logic to find the tenant from the database
		sqlDB, err := getDatabaseInstance()
		if err != nil {
			return err
		}

		query := "SELECT * FROM " + tenantModel.TableName() + " WHERE " + tenantModel.GetTenantKeyName() + " = ?"
		if err := sqlDB.QueryRow(query, tenant).Scan(&tenantModel.ID, &tenantModel.CreatedAt, &tenantModel.UpdatedAt, &tenantModel.DeletedAt, &tenantModel.Data); err != nil {
			return err
		}
		tenantInstance = tenantModel
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
	t.initialized = true

	// Emit events here if necessary

	return nil
}

func (t *Tenancy) End() error {
	// Emit events here if necessary

	if !t.initialized {
		return nil
	}

	// Emit events here if necessary

	t.initialized = false
	t.Tenant = models.Tenant{}

	return nil
}

func (t *Tenancy) GetBootstrappers() []contracts.TenancyBootstrapper {
	return t.GetConfig("tenancy.bootstrappers").([]contracts.TenancyBootstrapper)
}

func (t *Tenancy) RunForMultiple(tenants []interface{}, callback func(tenant models.Tenant) error) error {
	originalTenant := t.Tenant

	for _, tenant := range tenants {
		err := t.Initialize(tenant)
		if err != nil {
			return err
		}

		err = callback(t.Tenant)
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
