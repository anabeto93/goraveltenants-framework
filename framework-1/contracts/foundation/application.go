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
	Register(provider interface{}, force *bool) ServiceProvider
	RegisterDeferredProvider(provider string, service ...string)
	ResolveProvider(provider string) ServiceProvider
	Boot()
	Booting(callback func())
	Booted(callback func())
	BootstrapWith(bootstrappers []interface{})
	GetLocale() string
	GetNamespace() (string, error)
	GetProviders(provider interface{}) []ServiceProvider
	HasBeenBootstrapped() bool
	LoadDeferredProviders()
	SetLocale(locale string)
	ShouldSkipMiddleware() bool
	Terminating(callback interface{})
	Terminate()
	GetContainerInstance() container.Container
	IsProduction() bool
}
