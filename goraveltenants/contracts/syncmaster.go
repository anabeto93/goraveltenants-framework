package contracts

import (
	"gorm.io/gorm"
)

type SyncMaster interface {
	Syncable
	Tenants() *gorm.DB
	GetTenantModelName() string
}
