package goraveltenants

import (
	"github.com/goravel/framework/facades"
)

type CacheManager struct {
	store facades.Cache
}

func NewCacheManager(store facades.Cache) *CacheManager {
	return &CacheManager{
		store: store,
	}
}

func (m *CacheManager) Store(storeName ...string) facades.Cache {
	if len(storeName) > 0 {
		return m.store(storeName[0]).WithTag(m.getTag())
	}

	return m.store.WithTag(m.getTag())
}

func (m *CacheManager) Tags(names ...string) facades.Cache {
	return m.store.Tags(append(names, m.getTag())...)
}

func (m *CacheManager) getTag() string {
	return "tenant:" + tenant().GetTenantKey()
}
