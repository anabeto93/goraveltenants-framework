package contracts

import (
	"github.com/anabeto93/goraveltenants/contracts"
)

type DomainEvent struct {
	Domain Domain `json:"domain"`
}

func NewDomainEvent(domain Domain) *DomainEvent {
	return &DomainEvent{Domain: domain}
}
