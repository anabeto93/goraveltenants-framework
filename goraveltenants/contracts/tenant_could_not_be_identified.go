package contracts

import "fmt"

type TenantCouldNotBeIdentifiedException struct {
    message string
}

func NewTenantCouldNotBeIdentifiedException(message string) *TenantCouldNotBeIdentifiedException {
    return &TenantCouldNotBeIdentifiedException{
        message: message,
    }
}

func (e *TenantCouldNotBeIdentifiedException) Error() string {
    return fmt.Sprintf("Tenant could not be identified: %v", e.message)
}
