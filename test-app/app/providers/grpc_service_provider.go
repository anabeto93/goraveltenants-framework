package providers

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"

	"goravel/app/grpc"
	"goravel/routes"
)

type GrpcServiceProvider struct {
	*support.BaseServiceProvider
}

func (receiver *GrpcServiceProvider) Register() {
	//Add Grpc interceptors
	kernel := grpc.Kernel{}
	facades.Grpc.UnaryServerInterceptors(kernel.UnaryServerInterceptors())
	facades.Grpc.UnaryClientInterceptorGroups(kernel.UnaryClientInterceptorGroups())
}

func (receiver *GrpcServiceProvider) Boot() {
	//Add routes
	routes.Grpc()
}

func (receiver *GrpcServiceProvider) Name() string {
	return "GRPCServiceProvider"
}
