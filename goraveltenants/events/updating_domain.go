package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &UpdatingDomain{}

type UpdatingDomain struct {
	eventcontracts.BaseDomainEvent
}

func NewUpdatingDomainEvent(domain contracts.Domain) *UpdatingDomain {
	return eventcontracts.NewDomainEvent("UpdatingDomain", domain).(*UpdatingDomain)
}
