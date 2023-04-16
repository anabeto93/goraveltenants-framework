package events

import (
	"github.com/anabeto93/goraveltenants/contracts"
	eventcontracts "github.com/anabeto93/goraveltenants/events/contracts"
)

var _ eventcontracts.DomainEvent = &SavedDomain{}

type SavedDomain struct {
	eventcontracts.BaseDomainEvent
}

func NewSavedDomainEvent(domain contracts.Domain) *SavedDomain {
	return eventcontracts.NewDomainEvent("SavedDomain", domain).(*SavedDomain)
}
