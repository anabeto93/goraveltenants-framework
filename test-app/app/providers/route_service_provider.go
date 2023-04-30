package providers

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"

	"goravel/app/http"
	"goravel/routes"
)

type RouteServiceProvider struct {
	*support.BaseServiceProvider
}

func (receiver *RouteServiceProvider) Register() {
	//Add HTTP middlewares
	facades.Route.GlobalMiddleware(http.Kernel{}.Middleware()...)
}

func (receiver *RouteServiceProvider) Boot() {
	receiver.configureRateLimiting()

	routes.Web()
}

func (receiver *RouteServiceProvider) Name() string {
	return "RouteServiceProvider"
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
