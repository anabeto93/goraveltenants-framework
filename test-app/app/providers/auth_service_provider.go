package providers

import "github.com/goravel/framework/support"

type AuthServiceProvider struct {
	*support.BaseServiceProvider
}

func (receiver *AuthServiceProvider) Register() {

}

func (receiver *AuthServiceProvider) Boot() {
}

func (receiver *AuthServiceProvider) Name() string {
	return "providers.AuthServiceProvider"
}
