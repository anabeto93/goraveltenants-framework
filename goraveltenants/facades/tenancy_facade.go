package facades

import (
	"github.com/anabeto93/goraveltenants"
	"github.com/anabeto93/goraveltenants/contracts"
)

func init() {
	Tenancy = goraveltenants.NewTenancy()
}

var Tenancy contracts.Tenancy
