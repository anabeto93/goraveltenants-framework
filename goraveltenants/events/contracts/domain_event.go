package contracts

import (
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/contracts/event"
)

type DomainEvent interface {
	GetDomain() contracts.Domain
	Name() string
	Handle(args []event.Arg) ([]event.Arg, error)
}
type BaseDomainEvent struct {
	domain contracts.Domain
	name   string
}

func (bd *BaseDomainEvent) GetDomain() contracts.Domain {
	return bd.domain
}

func (bd *BaseDomainEvent) Name() string {
	return bd.name
}

func (bd *BaseDomainEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

func NewBaseTenant(name string, domain contracts.Domain) *BaseDomainEvent {
	return &BaseDomainEvent{name: name, domain: domain}
}

func NewDomainEvent(name string, domain contracts.Domain) DomainEvent {
	return NewBaseTenant(name, domain)
}
