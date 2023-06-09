package models

import (
	"strings"
	"time"

	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ contracts.Tenant = &Tenant{}
var _ contracts.HasInternalKeys = &Tenant{}

type Tenant struct {
	ID                uuid.UUID              `gorm:"type:uuid;primary_key;unique" json:"id"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         gorm.DeletedAt         `gorm:"index" json:"deleted_at"`
	TenancyDbUsername string                 `gorm:"size:500" json:"tenancy_db_username"`
	TenancyDbPassword string                 `gorm:"size:500" json:"tenancy_db_password"`
	KeyPath           string                 `json:"key_path"`
	Data              string                 `gorm:"type:json" json:"data"`
	attributes        map[string]interface{} `gorm:"-"`
}

func (t *Tenant) TableName() string {
	return "tenants"
}

func (t *Tenant) GetInternal(key string) interface{} {
	lowerKey := strings.ToLower(key)

	switch lowerKey {
	case "id":
		return t.ID
	case "created_at":
		return t.CreatedAt
	case "updated_at":
		return t.UpdatedAt
	case "deleted_at":
		return t.DeletedAt
	case "tenancy_db_username":
		return t.TenancyDbUsername
	case "tenancy_db_password":
		return t.TenancyDbPassword
	case "key_path":
		return t.KeyPath
	case "data":
		return t.Data
	}

	if val, ok := t.attributes[lowerKey]; ok {
		return val
	}

	return nil
}

func (t *Tenant) SetInternal(key string, value interface{}) {
	lowerKey := strings.ToLower(key)

	switch lowerKey {
	case "id":
		t.ID = value.(uuid.UUID)
	case "created_at":
		t.CreatedAt = value.(time.Time)
	case "updated_at":
		t.UpdatedAt = value.(time.Time)
	case "deleted_at":
		t.DeletedAt = value.(gorm.DeletedAt)
	case "tenancy_db_username":
		t.TenancyDbUsername = value.(string)
	case "tenancy_db_password":
		t.TenancyDbPassword = value.(string)
	case "key_path":
		t.KeyPath = value.(string)
	case "data":
		t.Data = value.(string)
	}

	t.attributes[lowerKey] = value
}

func (t *Tenant) GetAttributes() map[string]interface{} {
	var defaultAttributes = map[string]interface{}{
		"id":                  t.ID,
		"created_at":          t.CreatedAt,
		"updated_at":          t.UpdatedAt,
		"deleted_at":          t.DeletedAt,
		"tenancy_db_username": t.TenancyDbUsername,
		"tenancy_db_password": t.TenancyDbPassword,
		"key_path":            t.KeyPath,
		"data":                t.Data,
	}

	for k, val := range defaultAttributes {
		t.attributes[k] = val
	}

	return t.attributes
}

func (t *Tenant) GetTenantKeyName() string {
	return "id"
}

func (t *Tenant) GetTenantKey() interface{} {
	return t.GetInternal(t.GetTenantKeyName())
}

func (t *Tenant) Run(callback func(args ...interface{}) (interface{}, error)) (interface{}, error) {
	return callback(t)
}

func (t *Tenant) InternalPrefix() string {
	return "tenancy_"
}
