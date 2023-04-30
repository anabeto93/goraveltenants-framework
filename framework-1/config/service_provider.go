package config

import (
	"flag"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

type ServiceProvider struct {
	app foundationcontract.Application
}

func (config *ServiceProvider) Register() {
	var env *string
	if support.Env == support.EnvTest {
		testEnv := ".env"
		env = &testEnv
	} else {
		env = flag.String("env", ".env", "custom .env path")
		flag.Parse()
	}
	facades.Config = NewApplication(*env)
}

func (config *ServiceProvider) Boot() {

}

func (config *ServiceProvider) Name() string {
	return "ConfigServiceProvider"
}

func (config *ServiceProvider) CallBootingCallbacks() {
	//TODO implement me
}

func (config *ServiceProvider) CallBootedCallbacks() {
	//TODO implement me
}

func (config *ServiceProvider) NewInstance(application foundationcontract.Application) foundationcontract.ServiceProvider {
	config.app = application
	return config
}
