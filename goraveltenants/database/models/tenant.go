package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key;unique" json:"id"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	TenancyDbUsername      string         `gorm:"size:500" json:"tenancy_db_username"`
	TenancyDbPassword      string         `gorm:"size:500" json:"tenancy_db_password"`
	KeyPath                string         `json:"key_path"`
	Data                   string         `gorm:"type:json" json:"data"`
}

func (t *Tenant) TableName() string {
	return "tenants"
}