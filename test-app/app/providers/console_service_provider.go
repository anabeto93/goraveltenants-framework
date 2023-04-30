package providers

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"

	"goravel/app/console"
)

type ConsoleServiceProvider struct {
	*support.BaseServiceProvider
}

func (receiver *ConsoleServiceProvider) Register() {
	kernel := console.Kernel{}
	facades.Schedule.Register(kernel.Schedule())
	facades.Artisan.Register(kernel.Commands())
}

func (receiver *ConsoleServiceProvider) Boot() {

}

func (receiver *ConsoleServiceProvider) Name() string {
	return "ConsoleServiceProvider"
}
