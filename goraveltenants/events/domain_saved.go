package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &DomainSaved{}

type DomainSaved struct {
	eventcontracts.BaseDomainEvent
}

func NewDomainSavedEvent(domain contracts.Domain) *DomainSaved {
	return eventcontracts.NewDomainEvent("DomainSaved", domain).(*DomainSaved)
}
