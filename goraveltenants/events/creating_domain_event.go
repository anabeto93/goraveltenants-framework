package events

import (
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
	"github.com/anabeto93/goraveltenants/contracts"
)

var _ eventcontracts.DomainEvent = &CreatingDomain{}

type CreatingDomain struct {
	eventcontracts.BaseDomainEvent
}

func NewCreatingDomainEvent(domain contracts.Domain) *CreatingDomain {
	return eventcontracts.NewDomainEvent("CreatingDomain", domain).(*CreatingDomain)
}
