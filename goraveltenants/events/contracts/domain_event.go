package contracts

import (
	"github.com/anabeto93/goraveltenants/contracts"
)

type DomainEvent interface {
	GetDomain() contracts.Domain
	Name() string
}
type BaseDomainEvent struct {
	domain contracts.Domain `json:"domain"`
	name string	`json:"name"`
}

func (bd *BaseDomainEvent) GetDomain() contracts.Domain {
	return bd.domain
}

func (bd *BaseDomainEvent) Name() string {
	return bd.name
}

func NewBaseTenant(name string, domain contracts.Domain) *BaseDomainEvent {
	return &BaseDomainEvent{name: name, domain: domain}
}

func NewDomainEvent(name string, domain contracts.Domain) DomainEvent {
	return NewBaseTenant(name, domain)
}
