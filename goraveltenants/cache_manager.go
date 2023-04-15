package tenancy

import (
	"github.com/goravel/framework/cache"
)

type CacheManager struct {
	store cache.Store
}

func NewCacheManager(store cache.Store) *CacheManager {
	return &CacheManager{
		store: store,
	}
}

func (m *CacheManager) Store(storeName ...string) cache.Store {
	if len(storeName) > 0 {
		return m.store.Store(storeName[0]).WithTag(m.getTag())
	}

	return m.store.WithTag(m.getTag())
}

func (m *CacheManager) Tags(names ...string) cache.Store {
	return m.store.Tags(append(names, m.getTag())...)
}

func (m *CacheManager) getTag() string {
	return "tenant:" + tenant().GetTenantKey()
}
