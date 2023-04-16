package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &DomainDeleted{}

type DomainDeleted struct {
	eventcontracts.BaseDomainEvent
}

func NewDomainDeletedEvent(domain contracts.Domain) *DomainDeleted {
	return eventcontracts.NewDomainEvent("DomainDeleted", domain).(*DomainDeleted)
}
