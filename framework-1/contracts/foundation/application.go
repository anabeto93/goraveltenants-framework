package foundation

import "github.com/goravel/framework/contracts/container"

type Application interface {
	container.Container

	Version() string
	BasePath(path ...string) string
	BootstrapPath(path ...string) string
	ConfigPath(path ...string) string
	DatabasePath(path ...string) string
	ResourcePath(path ...string) string
	StoragePath(path ...string) string
	Environment(environments ...string) (string, bool)
	RunningInConsole() bool
	RunningUnitTests() bool
	MaintenanceMode() (MaintenanceMode, error)
	IsDownForMaintenance() (bool, error)
	RegisterConfiguredProviders()
	Register(provider interface{}, force *bool) (ServiceProvider, error)
	RegisterDeferredProvider(provider string, service *string) error
	ResolveProvider(provider string) (ServiceProvider, error)
	Boot()
	Booting(callback func(...interface{}))
	Booted(callback func(...interface{}))
	BootstrapWith(bootstrappers []interface{ Bootstrap(app Application) })
	GetLocale() string
	GetProviders(provider interface{}) []ServiceProvider
	HasBeenBootstrapped() bool
	LoadDeferredProviders() error
	SetLocale(locale string)
	ShouldSkipMiddleware() bool
	Terminating(callback interface{})
	Terminate()
	GetContainerInstance() container.Container
	IsProduction() bool
}
