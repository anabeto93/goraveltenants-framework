package models

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/anabeto93/goraveltenants/database"
)

var _ contracts.TenantWithDatabase = &BaseTenantWithDatabase{}

type BaseTenantWithDatabase struct {
	Tenant
	dbConfig contracts.DatabaseConfig
}

func NewBaseTenantWithDatabase(tenant *Tenant) *BaseTenantWithDatabase {
	dbConfig := database.NewDatabaseConfig(tenant)

	return &BaseTenantWithDatabase{
		Tenant:   *tenant,
		dbConfig: dbConfig,
	}
}

func (b *BaseTenantWithDatabase) Database() contracts.DatabaseConfig {
	return b.dbConfig
}
