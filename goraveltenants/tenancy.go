package goraveltenants

import (
	"context"
	"errors"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database/models"
)

var _ contracts.Tenancy = &Tenancy{}

type Tenancy struct {
	models.Tenant
	initialized       bool
	getBootstrappers  func(tenant contracts.Tenant) []contracts.TenancyBootstrapper
	withTransactionFn func(ctx context.Context, fn func(tx orm.Transaction) error) error
}

func NewTenancy(withTransactionFn func(ctx context.Context, fn func(tx orm.Transaction) error) error) *Tenancy {
	return &Tenancy{
		Tenant:            nil,
		initialized:       false,
		getBootstrappers:  nil,
		withTransactionFn: withTransactionFn,
	}
}

func (t *Tenancy) Initialize(tenant contracts.Tenant) error {
	if tenant == nil {
		return errors.New("tenant cannot be nil")
	}

	if t.initialized && t.Tenant.GetTenantKey() == tenant.GetTenantKey() {
		return nil
	}

	if t.initialized {
		t.End()
	}

	t.Tenant = tenant
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
	t.Tenant = nil

	return nil
}

func (t *Tenancy) GetBootstrappers() []contracts.TenancyBootstrapper {
	resolve := t.getBootstrappers
	if resolve == nil {
		resolve = func(tenant contracts.Tenant) []contracts.TenancyBootstrapper {
			bootstrappers := facades.Config.Get("tenancy.bootstrappers").([]contracts.TenancyBootstrapper)
			return bootstrappers
		}
	}

	return resolve(t.Tenant)
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
		t.Initialize(originalTenant)
	} else {
		t.End()
	}

	return nil
}

func (t *Tenancy) Query() orm.Query {
	return t.Tenant.Model().Query()
}

func (t *Tenancy) WithTransaction(ctx context.Context, fn func(tx orm.Transaction) error) error {
	return t.withTransactionFn(ctx, fn)
}
