package foundation

import (
	"fmt"
	"github.com/goravel/framework/container"
	containercontract "github.com/goravel/framework/contracts/container"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/filesystem"
	"os"
	"path/filepath"
	"reflect"
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

func (app *Application) IsLocal() bool {
	str := support.Env
	if str != "" {
		return str == "local"
	}
	env := facades.Config.GetString("app.env")
	return env == "local"
}

func (app *Application) IsProduction() bool {
	str := support.Env
	if str != "" {
		return str == "production"
	}

	env := facades.Config.GetString("app.env")
	return env == "production"
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

func (app *Application) Make(abstract string, parameters ...interface{}) (interface{}, error) {
	abstract = app.GetAlias(abstract)
	app.loadDeferredProviderIfNeeded(abstract)

	return app.Container.Make(abstract, parameters...)
}

func (app *Application) Resolve(abstract string, raiseEvents bool, parameters ...interface{}) (interface{}, error) {
	abstract = app.GetAlias(abstract)
	app.loadDeferredProviderIfNeeded(abstract)

	return app.Container.Resolve(abstract, raiseEvents, parameters...)
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

func (app *Application) GetContainerInstance() containercontract.Container {
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

func (app *Application) RegisterConfiguredProviders() error {
	manifest, err := app.Make("PackageManifest")
	if err != nil {
		return fmt.Errorf("PackageManifest not found in the Application")
	}
	packageManifest := manifest.(*PackageManifest)
	providers, err := packageManifest.Providers()
	if err != nil {
		return fmt.Errorf("error getting providers from PackageManifest: %v", err)
	}
	aliases, err := packageManifest.Aliases()
	if err != nil {
		return fmt.Errorf("error getting aliases from PackageManifest: %v", err)
	}

	// Register the providers
	for _, provider := range providers {
		force := false
		_, err := app.Register(provider, &force)
		if err != nil {
			return fmt.Errorf("error registering provider: %v", err)
		}
	}

	for alias, provider := range aliases {
		err := app.Alias(alias, provider.Name())
		if err != nil {
			return fmt.Errorf("error registering alias: %v", err)
		}
	}

	return nil
}

func (app *Application) Terminating(callback func(...interface{})) {
	callbacks := app.terminatingCallbacks
	callbacks = append(callbacks, callback)
	app.terminatingCallbacks = callbacks
}

func (app *Application) Terminate() {
	for _, callback := range app.terminatingCallbacks {
		_, _ = app.Call(callback)
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

func (app *Application) Bound(abstract string) bool {
	isDeferred := app.isDeferredService(abstract)
	isBound := app.Container.Bound(abstract)

	return isDeferred || isBound
}

func (app *Application) isBooted() bool {
	return app.booted
}

func (app *Application) GetLocale() string {
	return facades.Config.GetString("app.locale")
}

func (app *Application) CurrentLocale() string {
	return app.GetLocale()
}

func (app *Application) GetFallbackLocale() string {
	return facades.Config.GetString("app.fallback_locale")
}

func (app *Application) SetLocale(locale string) {
	facades.Config.Add("app.locale", locale)
	// TODO: set translator locale
	// TODO: dispatch LocaleUpdated event
}

func (app *Application) SetFallbackLocale(fallbackLocal string) {
	facades.Config.Add("app.fallback_locale", fallbackLocal)
	// TODO: set translator fallback locale
}

func (app *Application) IsLocale(locale string) bool {
	return app.GetLocale() == locale
}

// ShouldSkipMiddleware Determine if middleware has been disabled for the application
func (app *Application) ShouldSkipMiddleware() bool {
	isBound := app.Bound("middleware.disable")
	if !isBound {
		return false
	}

	disabled, err := app.Make("middleware.disable")
	if err != nil {
		return false
	}
	disabledMiddleware := disabled.(bool)
	return isBound && disabledMiddleware == true
}

func (app *Application) loadDeferredProviderIfNeeded(abstract string) {
	isDeferred := app.isDeferredService(abstract)
	notAnInstance := !app.IsInstance(abstract)

	if isDeferred && notAnInstance {
		_ = app.loadDeferredProvider(abstract)
	}
}

func (app *Application) isDeferredService(abstract string) bool {
	_, ok := app.deferredServices[abstract]
	return ok
}

func (app *Application) LoadDeferredProviders() error {
	currentServices := app.deferredServices

	for name, _ := range currentServices {
		if err := app.loadDeferredProvider(name); err != nil {
			return err
		}
	}

	app.deferredServices = make(map[string]contracts.ServiceProvider)
	return nil
}

func (app *Application) loadDeferredProvider(abstract string) error {
	if !app.isDeferredService(abstract) {
		return nil
	}

	provider := app.deferredServices[abstract]
	var isLoaded bool
	if loaded, found := app.loadedProviders[provider.Name()]; !found {
		isLoaded = false
	} else {
		isLoaded = loaded
	}

	if !isLoaded {
		return app.RegisterDeferredProvider(provider, &abstract)
	}
	return nil
}

func (app *Application) RegisterDeferredProvider(provider contracts.ServiceProvider, service *string) error {
	if service != nil {
		if _, ok := app.deferredServices[*service]; ok {
			delete(app.deferredServices, *service)
		}
	}

	instance := provider.NewInstance(app)
	notForced := false
	tempInstance, err := app.Register(instance, &notForced)
	if err != nil {
		return err
	}
	instance = tempInstance

	if !app.isBooted() {
		app.Booting(func(args ...interface{}) {
			app.bootProvider(instance)
		})
	}
	return nil
}

func (app *Application) Booting(callback func(params ...interface{})) {
	app.bootingCallbacks = append(app.bootedCallbacks, callback)
}

func (app *Application) Booted(callback func(params ...interface{})) {
	app.bootedCallbacks = append(app.bootedCallbacks, callback)

	if app.isBooted() {
		callback(app)
	}
}

func (app *Application) fireAppCallbacks(callbacks []func(...interface{})) {
	for _, callback := range callbacks {
		callback(app)
	}
}

func (app *Application) Register(provider interface{}, force *bool) (contracts.ServiceProvider, error) {
	if registered := app.getProvider(provider); registered != nil && *force == false {
		return registered, nil
	}
	var providerFound contracts.ServiceProvider
	if providerName, ok := provider.(string); ok {
		tempProvider, err := app.ResolveProvider(providerName)
		if err != nil {
			return nil, err
		}
		providerFound = tempProvider
	} else {
		providerFound = provider.(contracts.ServiceProvider)
	}

	// If there are bindings / singletons set as properties on the provider we
	// will spin through them and register them with the application, which
	// serves as a convenience layer while registering a lot of bindings.
	providerType := reflect.TypeOf(providerFound)

	bindingsField, bindingsExist := providerType.FieldByName("bindings")
	singletonsField, singletonsExist := providerType.FieldByName("singletons")

	if bindingsExist {
		value := reflect.ValueOf(providerFound).FieldByName(bindingsField.Name)
		bindings := value.Interface().(map[string]interface{})

		for key, val := range bindings {
			if err := app.Bind(key, val, false); err != nil {
				return nil, err
			}
		}
	}

	if singletonsExist {
		value := reflect.ValueOf(providerFound).FieldByName(singletonsField.Name)
		singletons := value.Interface().(map[string]interface{})

		for key, val := range singletons {
			if err := app.Singleton(key, val); err != nil {
				return nil, err
			}
		}
	}

	app.markAsRegistered(providerFound)

	// If the application has already booted, we will call this boot method on
	// the provider class so it has an opportunity to do its boot logic and
	// will be ready for any usage by this developer's application logic.
	if app.isBooted() {
		app.bootProvider(providerFound)
	}

	return providerFound, nil
}

// getProvider Get the registered service provider instance if it exists
func (app *Application) getProvider(provider interface{}) contracts.ServiceProvider {
	providers := app.GetProviders(provider)
	if len(providers) > 0 {
		return providers[0]
	}
	return nil
}

// GetProviders Get the registered service provider instances if any exist
func (app *Application) GetProviders(provider interface{}) []contracts.ServiceProvider {
	var name string
	switch provider.(type) {
	case string:
		name = provider.(string)
	case contracts.ServiceProvider:
		p := provider.(contracts.ServiceProvider)
		name = p.Name()
	}
	var providers []contracts.ServiceProvider

	for _, pvd := range app.serviceProviders {
		if pvd.Name() == name {
			providers = append(providers, pvd)
		}
	}

	return providers
}

func (app *Application) bootProvider(provider contracts.ServiceProvider) {
	provider.CallBootingCallbacks()

	provider.Boot()

	provider.CallBootedCallbacks()
}

func (app *Application) ResolveProvider(provider string) (contracts.ServiceProvider, error) {
	providerType := reflect.TypeOf(provider)
	if providerType == nil {
		return nil, fmt.Errorf("could not find provider with name: %s", provider)
	}

	providerValue := reflect.New(providerType)
	providerSrv, ok := providerValue.Interface().(contracts.ServiceProvider)
	if !ok {
		return nil, fmt.Errorf("provider %s does not implement ServiceProvider", provider)
	}

	return providerSrv.NewInstance(app), nil
}

func (app *Application) markAsRegistered(provider contracts.ServiceProvider) {
	currentProviders := app.serviceProviders
	currentProviders = append(currentProviders, provider)

	app.serviceProviders = currentProviders
	app.loadedProviders[provider.Name()] = true
}

// Boot the application's service providers
func (app *Application) Boot() {
	if app.isBooted() {
		return
	}

	// Once the application has booted we will also fire some "booted" callbacks
	// for any listeners that need to do work after this initial booting gets
	// finished. This is useful when ordering the boot-up processes we run.
	app.fireAppCallbacks(app.bootingCallbacks)

	for _, provider := range app.serviceProviders {
		app.bootProvider(provider)
	}

	app.booted = true

	app.fireAppCallbacks(app.bootedCallbacks)
}

func (app *Application) Flush() error {
	if err := app.Container.Flush(); err != nil {
		return err
	}

	app.loadedProviders = make(map[string]bool)
	app.bootedCallbacks = []func(...interface{}){}
	app.bootingCallbacks = []func(...interface{}){}
	app.deferredServices = make(map[string]contracts.ServiceProvider)
	app.serviceProviders = []contracts.ServiceProvider{}
	app.SetResolvingCallbacks(make(map[string][]func(...interface{}) interface{}))
	app.terminatingCallbacks = []func(...interface{}){}
	app.SetBeforeResolvingCallbacks(make(map[string][]func(...interface{}) interface{}))
	app.SetAfterResolvingCallbacks(make(map[string][]func(...interface{}) interface{}))
	app.SetGlobalBeforeResolvingCallbacks([]func(...interface{}) interface{}{})
	app.SetGlobalResolvingCallbacks([]func(...interface{}) interface{}{})
	app.SetGlobalAfterResolvingCallbacks([]func(...interface{}) interface{}{})

	return nil
}

func (app *Application) BootstrapWith(bootstrappers []interface{ Bootstrap(app foundationcontract.Application); Name() string }) {
	app.hasBeenBootstrapped = true

	for _, bootstrapper := range bootstrappers {
		app.dispatchEvent(fmt.Sprintf("bootstrapping: %s", bootstrapper.Name()))
		bootstrapper.Bootstrap(app)
		app.dispatchEvent(fmt.Sprintf("bootstrapped: %s", bootstrapper.Name()))
	}
}

func (app *Application) HasBeenBootstrapped() bool {
	return app.hasBeenBootstrapped
}

func (app *Application) dispatchEvent(eventName string) {

}

// Boot Register and bootstrap configured service providers.
func (app *Application) Bootx() {
	app.RegisterConfiguredServiceProviders()
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

// RegisterConfiguredServiceProviders Register configured service providers.
func (app *Application) RegisterConfiguredServiceProviders() {
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
