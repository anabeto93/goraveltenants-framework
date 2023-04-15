package exceptions

import (
    "fmt"
    "github.com/anabeto93/goraveltenants/contracts"
)

type TenantCouldNotBeIdentifiedOnDomainException struct {
    domain string
}

func NewTenantCouldNotBeIdentifiedOnDomainException(domain string) *TenantCouldNotBeIdentifiedOnDomainException {
    return &TenantCouldNotBeIdentifiedOnDomainException{domain: domain}
}

func (e *TenantCouldNotBeIdentifiedOnDomainException) Error() string {
    return fmt.Sprintf("Tenant could not be identified on domain %s", e.domain)
}

func (e *TenantCouldNotBeIdentifiedOnDomainException) Unwrap() error {
    return contracts.NewTenantCouldNotBeIdentifiedException(e.Error())
}
