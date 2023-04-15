package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantPivot struct {
	gorm.Model
	ID        uint      `gorm:"primary_key" json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ParentID  uuid.UUID `json:"parent_id"`
}

func (p *TenantPivot) TableName() string {
	return "tenant_pivots"
}

// BeforeSave is a GORM hook that triggers a sync event before saving the TenantPivot.
func (p *TenantPivot) BeforeSave(tx *gorm.DB) (err error) {
	parent := &Tenant{} // Or whatever the parent model is
	result := tx.Where("id = ?", p.ParentID).First(parent)
	if result.Error != nil {
		return result.Error
	}

	if syncable, ok := parent.(Syncable); ok {
		syncable.TriggerSyncEvent()
	}

	return nil
}
