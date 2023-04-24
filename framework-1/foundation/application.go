package foundation

import (
	"github.com/goravel/framework/container"
	containercontract "github.com/goravel/framework/contracts/container"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/filesystem"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/config"
	"github.com/goravel/framework/contracts"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

var _ containercontract.Container = &Application{}
var _ foundationcontract.Application = &Application{}

const VERSION = "1.10.0" // Bowen will correct this later

/*func init() {
	setEnv()

	app := Application{}
	app.RegisterBaseServiceProviders()
	app.bootBaseServiceProviders()
}*/

type Application struct {
	*container.Container
	basePath             string
	hasBeenBootstrapped  bool
	booted               bool
	bootingCallbacks     []func(...interface{})
	bootedCallbacks      []func(...interface{})
	terminatingCallbacks []func(...interface{})
	serviceProviders     []contracts.ServiceProvider
	loadedProviders      map[string]bool
	deferredServices     map[string]contracts.ServiceProvider
	appPath              string
	databasePath         string
	langPath             string
	storagePath          string
	environmentPath      string
	environmentFile      string
	isRunningInConsole   *bool
}

func NewApplication(basePath string) *Application {
	app := &Application{
		Container:            container.NewContainer(),
		basePath:             "",
		hasBeenBootstrapped:  false,
		booted:               false,
		bootingCallbacks:     []func(...interface{}){},
		bootedCallbacks:      []func(...interface{}){},
		terminatingCallbacks: []func(...interface{}){},
		serviceProviders:     []contracts.ServiceProvider{},
		loadedProviders:      make(map[string]bool),
		deferredServices:     make(map[string]contracts.ServiceProvider),
		appPath:              "",
		databasePath:         "",
		langPath:             "",
		storagePath:          "",
		environmentPath:      "",
		environmentFile:      "",
		isRunningInConsole:   nil,
	}

	if basePath != "" {
		app.basePath = basePath
	}

	app.registerBaseBindings()
	app.RegisterBaseServiceProviders()
	app.registerCoreContainerAliases()

	return app
}

func (app *Application) Version() string {
	return VERSION
}

func (app *Application) BasePath(path ...string) string {
	if len(path) > 0 {
		return filepath.Join(append([]string{app.basePath}, path...)...)
	}
	return app.basePath
}

func (app *Application) BootstrapPath(path ...string) string {
	bootstrapPath := filepath.Join(app.basePath, "bootstrap")
	if len(path) > 0 {
		return filepath.Join(append([]string{bootstrapPath}, path...)...)
	}
	return bootstrapPath
}

func (app *Application) ConfigPath(path ...string) string {
	configPath := filepath.Join(app.basePath, "config")
	if len(path) > 0 {
		return filepath.Join(append([]string{configPath}, path...)...)
	}
	return configPath
}

func (app *Application) DatabasePath(path ...string) string {
	if app.databasePath != "" {
		if len(path) > 0 {
			return filepath.Join(append([]string{app.databasePath}, path...)...)
		}
		return app.databasePath
	}

	databasePath := filepath.Join(app.basePath, "database")
	if len(path) > 0 {
		return filepath.Join(append([]string{databasePath}, path...)...)
	}
	return databasePath
}

func (app *Application) ResourcePath(path ...string) string {
	resourcePath := filepath.Join(app.basePath, "resources")
	if len(path) > 0 {
		return filepath.Join(append([]string{resourcePath}, path...)...)
	}
	return resourcePath
}

func (app *Application) StoragePath(path ...string) string {
	storagePath := filepath.Join(app.basePath, "storage")
	if len(path) > 0 {
		return filepath.Join(append([]string{storagePath}, path...)...)
	}
	return storagePath
}

func (app *Application) Environment(environments ...string) (string, bool) {
	if len(environments) == 0 {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "production"
		}
		return env, true
	}

	currentEnv, _ := app.Environment()
	for _, env := range environments {
		if env == currentEnv {
			return env, true
		}
	}

	return "", false
}

func (app *Application) RunningInConsole() bool {
	if app.isRunningInConsole != nil {
		return *app.isRunningInConsole
	}

	app.isRunningInConsole = new(bool) // because it is a pointer
	*app.isRunningInConsole = os.Getenv("APP_RUNNING_IN_CONSOLE") == "1"
	return *app.isRunningInConsole
}

func (app *Application) RunningUnitTests() bool {
	if app.Bound("env") {
		env, _ := app.Resolve("env", false)
		envStr, ok := env.(string)
		if ok {
			return envStr == "testing"
		}
	}
	return false
}

func (app *Application) MaintenanceMode() (foundationcontract.MaintenanceMode, error) {
	instance, err := app.Make("MaintenanceModeContract")
	if err != nil {
		return nil, err
	}

	return instance.(foundationcontract.MaintenanceMode), nil
}

func (app *Application) IsDownForMaintenance() (bool, error) {
	mode, err := app.MaintenanceMode()
	if err != nil {
		return false, err
	}

	return mode.Active(), nil
}

func (app *Application) GetContainerInstance() *container.Container {
	return app.Container
}

func (app *Application) registerBaseBindings() {
	_, _ = app.Instance("app", app)
	_, _ = app.Instance("Container", app)
	// Mix Later

	// In order to be able to use facades.Config in filesystem.NewStorage(), facades.Config has to first be registered
	fCon := &config.ServiceProvider{}
	fCon.Register()
	// Now Register the PackageManifest
	_ = app.Singleton("PackageManifest", func() *PackageManifest {
		return NewPackageManifest(filesystem.NewStorage(), app.basePath, "vendor/")
	})
}

func (app *Application) registerCoreContainerAliases() {
	coreAliases := map[string][]string{
		"app": {"Container", "Application"},
	}

	for key, aliases := range coreAliases {
		for _, alias := range aliases {
			_ = app.Alias(key, alias)
		}
	}
}

func (app *Application) getLoadedProviders() map[string]bool {
	return app.loadedProviders
}

func (app *Application) providerIsLoaded(provider string) bool {
	loaded, ok := app.loadedProviders[provider]
	if !ok {
		return false
	}
	return loaded
}

func (app *Application) getDeferredServices() map[string]contracts.ServiceProvider {
	return app.deferredServices
}

func (app *Application) setDeferredServices(services map[string]contracts.ServiceProvider) {
	app.deferredServices = services
}

func (app *Application) addDeferredServices(services map[string]contracts.ServiceProvider) {
	currentServices := app.deferredServices
	for name, provider := range services {
		currentServices[name] = provider
	}
	app.deferredServices = currentServices
}

func (app *Application) isBooted() bool {
	return app.booted
}

func (app *Application) loadDeferredProviderIfNeeded(abstract string) {
	isDeferred := app.isDeferredService(abstract)
	notAnInstance := !app.IsInstance(abstract)

	if isDeferred && notAnInstance {
		app.loadDeferredProvider(abstract)
	}
}

func (app *Application) isDeferredService(abstract string) bool {
	_, ok := app.deferredServices[abstract]
	return ok
}

func (app *Application) LoadDeferredProviders() {
	currentServices := app.deferredServices

	for name, _ := range currentServices {
		app.loadDeferredProvider(name)
	}

	app.deferredServices = make(map[string]contracts.ServiceProvider)
}

func (app *Application) loadDeferredProvider(abstract string) {
	if !app.isDeferredService(abstract) {
		return
	}

	provider := app.deferredServices[abstract]
	var isLoaded bool
	_, loaded := app.loadedProviders[]
}

func (app *Application) bootProvider(provider contracts.ServiceProvider) {
	provider.CallBootingCallbacks()

	provider.Boot()

	provider.CallBootedCallbacks()
}

func (app *Application) markAsRegistered(provider contracts.ServiceProvider) {
	currentProviders := app.serviceProviders
	currentProviders = append(currentProviders, provider)

	app.serviceProviders = currentProviders
	app.loadedProviders[provider.Name()] = true
}

// Boot Register and bootstrap configured service providers.
func (app *Application) Bootx() {
	app.registerConfiguredServiceProviders()
	app.bootConfiguredServiceProviders()

	app.bootArtisan()
	setRootPath()
}

// bootArtisan Boot artisan command.
func (app *Application) bootArtisan() {
	facades.Artisan.Run(os.Args, true)
}

// getBaseServiceProviders Get base service providers.
func (app *Application) getBaseServiceProviders() []contracts.ServiceProvider {
	return []contracts.ServiceProvider{
		&config.ServiceProvider{},
	}
}

// getConfiguredServiceProviders Get configured service providers.
func (app *Application) getConfiguredServiceProviders() []contracts.ServiceProvider {
	return facades.Config.Get("app.providers").([]contracts.ServiceProvider)
}

// RegisterBaseServiceProviders Register base service providers.
func (app *Application) RegisterBaseServiceProviders() {
	app.registerServiceProviders(app.getBaseServiceProviders())
}

// bootBaseServiceProviders Bootstrap base service providers.
func (app *Application) bootBaseServiceProviders() {
	app.bootServiceProviders(app.getBaseServiceProviders())
}

// registerConfiguredServiceProviders Register configured service providers.
func (app *Application) registerConfiguredServiceProviders() {
	app.registerServiceProviders(app.getConfiguredServiceProviders())
}

// bootConfiguredServiceProviders Bootstrap configured service providers.
func (app *Application) bootConfiguredServiceProviders() {
	app.bootServiceProviders(app.getConfiguredServiceProviders())
}

// registerServiceProviders Register service providers.
func (app *Application) registerServiceProviders(serviceProviders []contracts.ServiceProvider) {
	for _, serviceProvider := range serviceProviders {
		serviceProvider.Register()
	}
}

// bootServiceProviders Bootstrap service providers.
func (app *Application) bootServiceProviders(serviceProviders []contracts.ServiceProvider) {
	for _, serviceProvider := range serviceProviders {
		serviceProvider.Boot()
	}
}

func setEnv() {
	args := os.Args
	if strings.HasSuffix(os.Args[0], ".test") {
		support.Env = support.EnvTest
	}
	if len(args) >= 2 {
		if args[1] == "artisan" {
			support.Env = support.EnvArtisan
		}
	}
}

func setRootPath() {
	rootPath := getCurrentAbPath()

	// Hack air path
	airPath := "/storage/temp"
	if strings.HasSuffix(rootPath, airPath) {
		rootPath = strings.ReplaceAll(rootPath, airPath, "")
	}

	support.RootPath = rootPath
}
