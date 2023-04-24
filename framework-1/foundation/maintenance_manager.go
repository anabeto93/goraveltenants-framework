package foundation

import (
	"github.com/goravel/framework/contracts/cache"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
)

var _ foundationcontract.MaintenanceMode = &MaintenanceModeManager{}

type MaintenanceModeManager struct {
	*BaseManager
}

func (mm *MaintenanceModeManager) Activate(payload map[string]interface{}) {
	_, err := mm.Call("Activate", payload)
	if err != nil {
		panic(err)
	}
}

func (mm *MaintenanceModeManager) Deactivate() {
	_, err := mm.Call("Deactivate")
	if err != nil {
		panic(err)
	}
}

func (mm *MaintenanceModeManager) Active() bool {
	result, err := mm.Call("Active")
	if err != nil {
		return false
	}
	return result.(bool)
}

func (mm *MaintenanceModeManager) Data() map[string]interface{} {
	result, err := mm.Call("Data")
	if err != nil {
		return nil
	}
	return result.(map[string]interface{})
}

func (mm *MaintenanceModeManager) CreateFileDriver() *FileBasedMaintenanceMode {
	return NewFileBasedMaintenanceMode(mm.GetContainer())
}

func (mm *MaintenanceModeManager) CreateCacheDriver() *CacheBasedMaintenanceMode {
	cacheManager, err := mm.GetContainer().Make("cache")
	if err != nil {
		panic(err)
	}
	store := mm.config.Get("app.maintenance.store").(string)
	if store == "" {
		store = mm.config.GetString("cache.default")
	}
	return NewCacheBasedMaintenanceMode(cacheManager.(cache.Cache), store, "goravel:foundation:down")
}

func (mm *MaintenanceModeManager) GetDefaultDriver() string {
	return mm.config.GetString("app.maintenance.store", "file")
}
