package foundation

import (
	"fmt"
	"github.com/goravel/framework/contracts/cache"
)

type CacheBasedMaintenanceMode struct {
	cache cache.Cache
	store string
	key   string
}

func NewCacheBasedMaintenanceMode(cache cache.Cache, store, key string) *CacheBasedMaintenanceMode {
	return &CacheBasedMaintenanceMode{
		cache: cache,
		store: store,
		key:   key,
	}
}

func (c *CacheBasedMaintenanceMode) Activate(payload map[string]interface{}) error {
	if ok := c.getStore().Forever(c.key, payload); !ok {
		return fmt.Errorf("failed to Activate")
	}
	return nil
}

func (c *CacheBasedMaintenanceMode) Deactivate() error {
	if ok := c.getStore().Forget(c.key); !ok {
		return fmt.Errorf("failed to Deactivate")
	}
	return nil
}

func (c *CacheBasedMaintenanceMode) Active() bool {
	return c.getStore().Has(c.key)
}

func (c *CacheBasedMaintenanceMode) Data() (map[string]interface{}, error) {
	data := c.getStore().Get(c.key)
	if data == nil {
		return nil, fmt.Errorf("data not found")
	}
	return data.(map[string]interface{}), nil
}

func (c *CacheBasedMaintenanceMode) getStore() cache.Driver {
	return c.cache.Store(c.store)
}
