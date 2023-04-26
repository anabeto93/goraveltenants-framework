package foundation

type ServiceProvider interface {
	//Boot any application services after register.
	Boot()
	//Register any application services.
	Register()
	//Name required to uniquely identify all service providers
	Name() string
	//CallBootingCallbacks calls the booting callbacks
	CallBootingCallbacks()
	//CallBootedCallbacks calls the booted callbacks
	CallBootedCallbacks()
	//NewInstance returns an instance of the current service provider
	NewInstance(application Application) ServiceProvider
}
