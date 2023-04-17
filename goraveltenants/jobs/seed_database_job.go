package jobs

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/facades"
)

var _ Job = &SeedDatabaseJob{}

type SeedDatabaseJob struct {
	Tenant contracts.TenantWithDatabase
}

func NewSeedDatabaseJob(tenant contracts.TenantWithDatabase) *SeedDatabaseJob {
	return &SeedDatabaseJob{
		Tenant: tenant,
	}
}

func (job *SeedDatabaseJob) Execute(args ...interface{}) error {
	facades.Artisan.Call("tenants:seed --tenants" + job.Tenant.GetTenantKey().(string))

	return nil
}
