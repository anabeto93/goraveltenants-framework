package container

type Container interface {
	Bound(abstract string) bool
	Alias(abstract string, alias string) error
	Tag(abstracts interface{}, tags ...interface{}) error
	Tagged(tag string) (map[string]interface{}, error)
	Bind(abstract string, concrete interface{}, shared bool) error
	BindIf(abstract string, concrete interface{}, shared bool) error
	Singleton(abstract string, concrete interface{}) error
	SingletonIf(abstract string, concrete interface{}) error
	Scoped(abstract string, concrete interface{}) error
	ScopedIf(abstract string, concrete interface{}) error
	Extend(abstract string, closure func(...interface{}) (interface{}, error)) error
	Instance(abstract string, instance interface{}) (interface{}, error)
	ForgetInstance(abstract string)
	AddContextualBinding(concrete string, abstract string, implementation interface{}) error
	When(concrete ...interface{}) ContextualBindingBuilder
	Factory(abstract string) func() (interface{}, error)
	Flush() error
	Make(abstract string, parameters ...interface{}) (interface{}, error)
	Call(callback interface{}, parameters ...interface{}) (interface{}, error)
	Resolve(abstract string, raiseEvents bool, parameters ...interface{}) (interface{}, error)
	Resolved(abstract string) bool
	BeforeResolving(abstract interface{}, callback *func(string, Container, ...interface{}) error) error
	Resolving(abstract interface{}, callback *func(string, Container, ...interface{}) error) error
	AfterResolving(abstract interface{}, callback *func(string, Container, ...interface{}) error) error
	IsInstance(abstract string) bool
}
