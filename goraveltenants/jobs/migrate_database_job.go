package jobs

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/facades"
)

var _ Job = &MigrateDatabaseJob{}

type MigrateDatabaseJob struct {
	Tenant contracts.TenantWithDatabase
}

func NewMigrateDatabaseJob(tenant contracts.TenantWithDatabase) *MigrateDatabaseJob {
	return &MigrateDatabaseJob{
		Tenant: tenant,
	}
}

func (job *MigrateDatabaseJob) Execute(args ...interface{}) error {
	facades.Artisan.Call("tenants:migrate --tenants" + job.Tenant.GetTenantKey().(string))

	return nil
}
