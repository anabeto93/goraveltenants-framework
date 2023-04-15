package contracts

type TenantCannotBeCreatedException interface {
	Reason() string
	error
}

type tenantCannotBeCreatedException struct {
	msg string
}

func (e *tenantCannotBeCreatedException) Error() string {
	return e.msg
}

func (e *tenantCannotBeCreatedException) Reason() string {
	return e.msg
}

func NewTenantCannotBeCreatedException(msg string) TenantCannotBeCreatedException {
	return &tenantCannotBeCreatedException{msg}
}
