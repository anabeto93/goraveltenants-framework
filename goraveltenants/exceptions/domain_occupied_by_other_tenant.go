package exceptions

type DomainOccupiedByOtherTenantException struct {
    Domain string
}

func (e DomainOccupiedByOtherTenantException) Error() string {
    return "The " + e.Domain + " domain is occupied by another tenant."
}
