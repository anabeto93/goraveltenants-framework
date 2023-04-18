package bootstrappers

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/foundation"
)

var _ contracts.TenancyBootstrapper = &CacheTenancyBootstrapper{}

type CacheTenancyBootstrapper struct {
	originalCache cache.Store
	app           *foundation.Application
}

func NewCacheTenancyBootstrapper(app *foundation.Application) *CacheTenancyBootstrapper {
	return &CacheTenancyBootstrapper{
		app: app,
	}
}

func (b *CacheTenancyBootstrapper) Bootstrap(tenant contracts.Tenant) {
	b.resetFacadeCache()

	b.originalCache = b.app.Cache()
	tenantCacheManager := cache.NewTenantCacheManager(b.app) // Adjust this line to create a tenant-specific cache.Store instance
	b.app.Bind("cache", tenantCacheManager)
}

func (b *CacheTenancyBootstrapper) Revert() {
	b.resetFacadeCache()

	if b.originalCache != nil {
		b.app.Bind("cache", b.originalCache)
		b.originalCache = nil
	}
}

func (b *CacheTenancyBootstrapper) resetFacadeCache() {
	cache.Flush()
}
