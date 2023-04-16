package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &DomainCreated{}

type DomainCreated struct {
	eventcontracts.BaseDomainEvent
}

func NewDomainCreatedEvent(domain contracts.Domain) *DomainCreated {
	return eventcontracts.NewDomainEvent("DomainCreated", domain).(*DomainCreated)
}
