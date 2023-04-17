package jobs

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

var _ Job = &DeleteDatabaseJob{}

type DeleteDatabaseJob struct {
	Tenant contracts.TenantWithDatabase
}

func NewDeleteDatabaseJob(tenant contracts.TenantWithDatabase) *DeleteDatabaseJob {
	return &DeleteDatabaseJob{
		Tenant: tenant,
	}
}

func (job *DeleteDatabaseJob) Execute(args ...interface{}) error {
	if err := facades.Event.Job(events.NewDeletingDatabase(job.Tenant), []event.Arg{}).Dispatch(); err != nil {
		return err
	}

	job.Tenant.Database().Manager().DeleteDatabase(job.Tenant)

	if err := facades.Event.Job(events.NewDatabaseDeleted(job.Tenant), []event.Arg{}).Dispatch(); err != nil {
		return err
	}

	return nil
}
