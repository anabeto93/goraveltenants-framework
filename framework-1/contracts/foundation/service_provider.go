package foundation

type ServiceProvider interface {
	//Boot any application services after register.
	Boot()
	//Register any application services.
	Register()
}
