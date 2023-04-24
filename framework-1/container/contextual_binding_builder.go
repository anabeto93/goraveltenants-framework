package container

import (
	configcontract "github.com/goravel/framework/contracts/config"
	containercontract "github.com/goravel/framework/contracts/container"
)

var _ containercontract.ContextualBindingBuilder = &ContextualBindingBuilder{}

type ContextualBindingBuilder struct {
	container containercontract.Container
	concrete  interface{}
	needs     string
}

func NewContextualBindingBuilder(container containercontract.Container, concrete ...interface{}) *ContextualBindingBuilder {
	return &ContextualBindingBuilder{
		container: container,
		concrete:  concrete,
	}
}

func (cbb *ContextualBindingBuilder) Needs(abstract string) containercontract.ContextualBindingBuilder {
	cbb.needs = abstract
	return cbb
}

func (cbb *ContextualBindingBuilder) Give(implementation interface{}) error {
	concretes := arrayWrap(cbb.concrete)

	for _, concrete := range concretes {
		c, _ := concrete.(string)
		if err := cbb.container.AddContextualBinding(c, cbb.needs, implementation); err != nil {
			return err
		}
	}

	return nil
}

func (cbb *ContextualBindingBuilder) GiveTagged(tag string) error {
	return cbb.Give(func(container containercontract.Container) []interface{} {
		taggedServices, _ := container.Tagged(tag)
		var services []interface{}
		for _, service := range taggedServices {
			services = append(services, service)
		}

		return services
	})
}

func (cbb *ContextualBindingBuilder) GiveConfig(key string, defaultValue interface{}) error {
	return cbb.Give(func(container containercontract.Container) interface{} {
		madeConfig, _ := container.Make("config")
		config := madeConfig.(configcontract.Config)
		return config.Get(key, defaultValue)
	})
}
