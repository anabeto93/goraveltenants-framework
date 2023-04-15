package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Domain struct {
	gorm.Model
	ID        uint      `gorm:"primary_key" json:"id"`
	Domain    string    `gorm:"size:255;unique" json:"domain"`
	TenantID  uuid.UUID `json:"tenant_id"`

	Tenant Tenant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tenant"`
}

func (d *Domain) TableName() string {
	return "domains"
}
