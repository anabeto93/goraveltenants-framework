package container

type ContextualBindingBuilder interface {
	Needs(abstract string) ContextualBindingBuilder
	Give(implementation interface{}) error
	GiveTagged(tag string) error
	GiveConfig(key string, defaultValue interface{}) error
}
