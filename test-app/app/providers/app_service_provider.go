package providers

import "github.com/goravel/framework/support"

type AppServiceProvider struct {
	*support.BaseServiceProvider
}

func (asp *AppServiceProvider) Register() {

}

func (asp *AppServiceProvider) Boot() {

}

func (asp *AppServiceProvider) Name() string {
	return "AppServiceProvider"
}
