package exceptions

import "fmt"

type NotASubdomainException struct {
	Hostname string
}

func (e *NotASubdomainException) Error() string {
	return fmt.Sprintf("Hostname %s does not include a subdomain.", e.Hostname)
}
