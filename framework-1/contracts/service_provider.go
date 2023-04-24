package contracts

import foundationcontract "github.com/goravel/framework/contracts/foundation"

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
	//NewInstance returns an instance of the service provider with app defined
	NewInstance(application foundationcontract.Application)
}
