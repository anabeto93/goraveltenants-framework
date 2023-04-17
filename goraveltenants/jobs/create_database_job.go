package jobs

import (
	"errors"
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database"
	"github.com/anabeto93/goraveltenants/events"
	"github.com/anabeto93/goraveltenants/exceptions"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

var _ Job = &CreateDatabaseJob{}

type CreateDatabaseJob struct {
	Tenant contracts.TenantWithDatabase
}

func NewCreateDatabaseJob(tenant contracts.TenantWithDatabase) *CreateDatabaseJob {
	return &CreateDatabaseJob{
		Tenant: tenant,
	}
}

func (job *CreateDatabaseJob) Execute(args ...interface{}) error {
	databaseManager, ok := args[0].(*database.DatabaseManager)
	if !ok {
		return errors.New("A *database.DatabaseManager is required in CreateDatabaseJob")
	}

	if err := facades.Event.Job(events.NewCreatingDatabase(job.Tenant), []event.Arg{}).Dispatch(); err != nil {
		return err
	}

	// Terminate execution if create_database is set to false
	if value, ok := job.Tenant.GetInternal("create_database").(bool); ok && !value {
		return nil
	}

	job.Tenant.Database().MakeCredentials()
	err := databaseManager.EnsureTenantCanBeCreated(job.Tenant)
	if err != nil {
		if _, ok := err.(*exceptions.TenantDatabaseAlreadyExistsException); ok {
			// Handle TenantDatabaseAlreadyExistsException if necessary
		} else if _, ok := err.(*exceptions.TenantDatabaseUserAlreadyExistsException); ok {
			// Handle TenantDatabaseUserAlreadyExistsException if necessary
		}
		return err
	}

	job.Tenant.Database().Manager().CreateDatabase(job.Tenant)

	// Ignoring event dispatch: DatabaseCreated
	if err := facades.Event.Job(events.NewDatabaseCreated(job.Tenant), []event.Arg{}).Dispatch(); err != nil {
		return err
	}

	return nil
}
