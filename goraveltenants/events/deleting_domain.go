package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &DeletingDomain{}

type DeletingDomain struct {
	eventcontracts.BaseDomainEvent
}

func NewDeletingDomainEvent(domain contracts.Domain) *DeletingDomain {
	return eventcontracts.NewDomainEvent("DeletingDomain", domain).(*DeletingDomain)
}
