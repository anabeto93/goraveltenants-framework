package exceptions

type TenancyNotInitializedException struct {
    message string
}

func (e *TenancyNotInitializedException) Error() string {
    if e.message != "" {
        return e.message
    }
    return "Tenancy is not initialized."
}
