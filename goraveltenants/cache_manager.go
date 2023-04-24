package goraveltenants

import (
	"errors"
	"fmt"
	"github.com/anabeto93/goraveltenants/facades"
	"github.com/goravel/framework/cache"
	frameworkfacades "github.com/goravel/framework/facades"
	"github.com/goravel/framework/foundation"
)

type TenantCacheManager struct {
	app          *foundation.Application
	tenantPrefix string
}

func NewTenantCacheManager(app *foundation.Application) *TenantCacheManager {
	return &TenantCacheManager{
		app:          app,
		tenantPrefix: "tenancy.cache.tag_base", // Replace this with your desired tenant prefix
	}
}

func (tcm *TenantCacheManager) Call(method string, parameters ...interface{}) (interface{}, error) {
	var tags []string

	tenant := facades.Tenancy.GetCurrentTenant() // Replace this with the function to get the current tenant

	if tenant == nil {
		return nil, errors.New("no active tenant")
	}

	tags = append(tags, fmt.Sprintf("%s%s", tcm.tenantPrefix, tenant.GetTenantKey()))

	if method == "Tags" {
		if len(parameters) != 1 {
			return nil, errors.New("Method Tags() takes exactly 1 argument")
		}

		names := parameters[0].([]string)
		tags = append(tags, names...)

		frameworkfacades.Cache.

		return tcm.app.Cache().Tags(tags), nil
	}

	return callWithTaggedCacheMethod(tcm.app.Cache().Tags(tags), method, parameters)
}

func callWithTaggedCacheMethod(taggedCache cache.TaggedCache, method string, parameters []interface{}) (interface{}, error) {
	// Implement reflection or type assertion to call the method on the taggedCache
	// You will need to handle each method individually based on the method name and parameters
}
