package goraveltenants

import "github.com/goravel/framework/contracts"

var _ contracts.ServiceProvider = &TenancyServiceProvider{}

type TenancyServiceProvider struct{}

func (tsp *TenancyServiceProvider) Register() {}

func (tsp *TenancyServiceProvider) Boot() {}
