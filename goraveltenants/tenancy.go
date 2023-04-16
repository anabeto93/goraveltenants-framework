package goraveltenants

import (
	"errors"
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database/models"
	"github.com/goravel/framework/facades"
)

var _ contracts.Tenancy = &Tenancy{}

type Tenancy struct {
	models.Tenant
	initialized      bool
	getBootstrappers func(tenant contracts.Tenant) []contracts.TenancyBootstrapper
}

func (t *Tenancy) Initialize(tenant contracts.Tenant) error {
	if tenant == nil {
		return errors.New("tenant cannot be nil")
	}

	if t.initialized && t.Tenant.GetTenantKey() == tenant.GetTenantKey() {
		return nil
	}

	if t.initialized {
		if err := t.End(); err != nil {
			return err
		}
	}

	temp := tenant.(models.Tenant)
	t.Tenant = tenant.(models.Tenant)
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
	return facades.Config.Get("tenancy.bootstrappers").([]contracts.TenancyBootstrapper)
}

func (t *Tenancy) RunForMultiple(tenants []contracts.Tenant, callback func(tenant contracts.Tenant) error) error {
	originalTenant := t.Tenant

	for _, tenant := range tenants {
		err := t.Initialize(tenant)
		if err != nil {
			return err
		}

		err = callback(tenant)
		if err != nil {
			return err
		}
	}

	if originalTenant != nil {
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
