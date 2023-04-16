package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &DomainUpdated{}

type DomainUpdated struct {
	eventcontracts.BaseDomainEvent
}

func NewDomainUpdatedEvent(domain contracts.Domain) *DomainUpdated {
	return eventcontracts.NewDomainEvent("DomainUpdated", domain).(*DomainUpdated)
}
