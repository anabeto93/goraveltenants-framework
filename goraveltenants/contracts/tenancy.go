package contracts

import (
	"context"
	"github.com/goravel/framework/contracts/database/orm"
)

type Tenancy interface {
	Tenant
	Initialize(tenant Tenant) error
	End() error
	GetBootstrappers() []TenancyBootstrapper
	RunForMultiple(tenants []Tenant, callback func(tenant Tenant) error) error
	Query() orm.Query
	WithTransaction(ctx context.Context, fn func(tx orm.Transaction) error) error
}
