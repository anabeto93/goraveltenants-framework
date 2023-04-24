package container

import (
	"fmt"
	containercontract "github.com/goravel/framework/contracts/container"
	"reflect"
	"sync"
)

var _ containercontract.Container = &Container{}

type binding struct {
	concrete interface{}
	shared   bool
}

type Container struct {
	bindings                       map[string]*binding
	instances                      map[string]interface{}
	aliases                        map[string]string
	tags                           map[string][]string
	reboundCallbacks               map[string][]func(container containercontract.Container, instance interface{})
	scopedInstances                []string
	abstractAliases                map[string][]interface{}
	contextualBindings             map[string]map[string]interface{}
	resolved                       map[string]bool
	globalBeforeResolvingCallbacks []func(...interface{}) interface{}
	beforeResolvingCallbacks       map[string][]func(...interface{}) interface{}
	extenders                      map[string][]func(...interface{}) (interface{}, error)
	resolvingCallbacks             map[string][]func(...interface{}) interface{}
	globalResolvingCallbacks       []func(...interface{}) interface{}
	afterResolvingCallbacks        map[string][]func(...interface{}) interface{}
	globalAfterResolvingCallbacks  []func(...interface{}) interface{}
	mu                             sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		bindings:                       make(map[string]*binding),
		instances:                      make(map[string]interface{}),
		aliases:                        make(map[string]string),
		tags:                           make(map[string][]string),
		reboundCallbacks:               make(map[string][]func(container containercontract.Container, instance interface{})),
		scopedInstances:                []string{},
		abstractAliases:                make(map[string][]interface{}),
		contextualBindings:             make(map[string]map[string]interface{}),
		resolved:                       make(map[string]bool),
		globalBeforeResolvingCallbacks: []func(...interface{}) interface{}{},
		beforeResolvingCallbacks:       make(map[string][]func(...interface{}) interface{}),
		extenders:                      make(map[string][]func(...interface{}) (interface{}, error)),
		resolvingCallbacks:             make(map[string][]func(...interface{}) interface{}),
		globalResolvingCallbacks:       []func(...interface{}) interface{}{},
		afterResolvingCallbacks:        make(map[string][]func(...interface{}) interface{}),
		globalAfterResolvingCallbacks:  []func(...interface{}) interface{}{},
	}
}

func (c *Container) Bound(abstract string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, bound := c.bindings[abstract]; bound {
		return true
	}

	_, aliasBound := c.aliases[abstract]
	return aliasBound
}

func (c *Container) Alias(abstract string, alias string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if abstract == alias {
		return fmt.Errorf("%s is aliased to itself", alias)
	}

	/*if !c.Bound(abstract) {// don't remember why I even put this here
		return containercontract.NewBindingResolutionError("Cannot alias, abstract type not bound.")
	}*/

	c.aliases[alias] = abstract
	abstractOnes := c.abstractAliases
	abstractOnes[abstract] = append(abstractOnes[abstract], alias)
	c.abstractAliases = abstractOnes
	return nil
}

func (c *Container) Tag(abstracts interface{}, tags ...interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	abstractsKeys, err := c.resolveKeys(abstracts)
	if err != nil {
		return err
	}

	tagKeys, err := c.resolveKeys(tags)
	if err != nil {
		return err
	}

	for _, tag := range tagKeys {
		if _, ok := c.tags[tag]; !ok {
			c.tags[tag] = make([]string, 0)
		}
		for _, abstract := range abstractsKeys {
			c.tags[tag] = append(c.tags[tag], abstract)
		}
	}

	return nil
}

func (c *Container) Tagged(tag string) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tagged := make(map[string]interface{})
	if taggedServices, ok := c.tags[tag]; ok {
		for _, abstract := range taggedServices {
			instance, err := c.Resolve(abstract, true)
			if err != nil {
				return nil, err
			}

			tagged[abstract] = instance
		}
	}

	return tagged, nil
}

func (c *Container) Bind(abstract string, concrete interface{}, shared bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.dropStaleInstances(abstract)

	if concrete == nil {
		concrete = abstract
	}

	if reflect.TypeOf(concrete).Kind() != reflect.Func {
		if _, ok := concrete.(string); !ok {
			return fmt.Errorf("the Bind() Argument #2 (concrete) must be of type string or function or nil")
		}
		concrete = c.getClosure(abstract, concrete)
	}

	c.bindings[abstract] = &binding{
		concrete: concrete,
		shared:   shared,
	}

	if _, resolved := c.instances[abstract]; resolved {
		c.rebound(abstract)
	}
	return nil
}

func (c *Container) BindIf(abstract string, concrete interface{}, shared bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, bound := c.bindings[abstract]; !bound {
		return c.Bind(abstract, concrete, shared)
	}
	return nil
}

func (c *Container) Singleton(abstract string, concrete interface{}) error {
	return c.Bind(abstract, concrete, true)
}

func (c *Container) SingletonIf(abstract string, concrete interface{}) error {
	return c.BindIf(abstract, concrete, true)
}

func (c *Container) Scoped(abstract string, concrete interface{}) error {
	return c.Bind(abstract, concrete, false)
}

func (c *Container) ScopedIf(abstract string, concrete interface{}) error {
	return c.BindIf(abstract, concrete, false)
}

func (c *Container) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.aliases = make(map[string]string)
	c.resolved = make(map[string]bool)
	c.bindings = make(map[string]*binding)
	c.instances = make(map[string]interface{})
	c.abstractAliases = make(map[string][]interface{})
	c.scopedInstances = []string{}

	return nil
}

func (c *Container) Extend(abstract string, closure func(...interface{}) (interface{}, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	abstract = c.getAlias(abstract)

	var err error
	if instance, ok := c.instances[abstract]; ok {
		c.instances[abstract], err = closure(instance, c)
		if err != nil {
			return err
		}
		c.rebound(abstract)
	} else {
		c.extenders[abstract] = append(c.extenders[abstract], closure)
		if c.Resolved(abstract) {
			c.rebound(abstract)
		}
	}

	return nil
}

func (c *Container) Resolved(abstract string) bool {
	if c.isAlias(abstract) {
		abstract = c.getAlias(abstract)
	}

	if _, ok := c.resolved[abstract]; ok {
		return true
	}

	_, resolved := c.instances[abstract]
	return resolved
}

func (c *Container) AddContextualBinding(concrete, abstract string, implementation interface{}) error {
	if c.contextualBindings == nil {
		c.contextualBindings = make(map[string]map[string]interface{})
	}

	if _, ok := c.contextualBindings[concrete]; !ok {
		c.contextualBindings[concrete] = make(map[string]interface{})
	}

	abstractAlias := c.getAlias(abstract)
	c.contextualBindings[concrete][abstractAlias] = implementation
	return nil
}

func (c *Container) When(concretes ...interface{}) containercontract.ContextualBindingBuilder {
	var aliases []interface{}

	for _, impl := range arrayWrap(concretes) {
		aliases = append(aliases, c.getAlias(impl.(string)))
	}
	builder := NewContextualBindingBuilder(c, aliases)
	return builder
}

func (c *Container) Factory(abstract string) func() (interface{}, error) {
	return func() (interface{}, error) {
		return c.Make(abstract)
	}
}

func (c *Container) MakeWith(abstract string, parameters []interface{}) (interface{}, error) {
	return c.Make(abstract, parameters...)
}

func (c *Container) Make(abstract string, parameters ...interface{}) (interface{}, error) {
	return c.Resolve(abstract, true, parameters...)
}

func (c *Container) Call(callback interface{}, parameters ...interface{}) (interface{}, error) {
	callableValue := reflect.ValueOf(callback)
	if callableValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("callback is not a function")
	}

	numArgs := callableValue.Type().NumIn()
	if numArgs != len(parameters) {
		return nil, fmt.Errorf("callback expects %d arguments, but %d were provided", numArgs, len(parameters))
	}

	argValues := make([]reflect.Value, len(parameters))
	for i, arg := range parameters {
		argValues[i] = reflect.ValueOf(arg)
	}

	result := callableValue.Call(argValues)
	if len(result) == 0 {
		return nil, nil
	}

	return result[0].Interface(), nil
}

func (c *Container) BeforeResolving(abstract interface{}, callback *func(string, containercontract.Container, ...interface{}) error) error {
	var key string

	switch v := abstract.(type) {
	case string:
		key = c.getAlias(v)
	case func(...interface{}) error:
		typedCallback := func(args ...interface{}) interface{} {
			currentCallback := abstract.(func(...interface{}) error)
			result := currentCallback(args...)
			return result
		}
		c.globalBeforeResolvingCallbacks = append(c.globalBeforeResolvingCallbacks, typedCallback)
		return nil
	case func(...interface{}) interface{}:
		typedCallback := abstract.(func(...interface{}) interface{})
		c.globalBeforeResolvingCallbacks = append(c.globalBeforeResolvingCallbacks, typedCallback)
		return nil
	default:
		return fmt.Errorf("abstract must be a string or a function")
	}

	if callback != nil {
		typedCallback := func(args ...interface{}) interface{} {
			deReferencedFunc := *callback
			result := deReferencedFunc(args[0].(string), args[1].(containercontract.Container), args[2:]...)
			return result
		}
		c.beforeResolvingCallbacks[key] = append(c.beforeResolvingCallbacks[key], typedCallback)
	}
	return nil
}

func (c *Container) Resolving(abstract interface{}, callback *func(string, containercontract.Container, ...interface{}) error) error {
	var key string

	switch v := abstract.(type) {
	case string:
		key = c.getAlias(v)
	case func(...interface{}) error:
		typedCallback := func(args ...interface{}) interface{} {
			currentCallback := abstract.(func(...interface{}) interface{})
			result := currentCallback(args...)
			return result
		}
		c.globalResolvingCallbacks = append(c.globalResolvingCallbacks, typedCallback)
		return nil
	case func(...interface{}) interface{}:
		typedCallback := abstract.(func(...interface{}) interface{})
		c.globalResolvingCallbacks = append(c.globalResolvingCallbacks, typedCallback)
	default:
		return fmt.Errorf("abstract must be a string or a function")
	}

	if callback != nil {
		typedCallback := func(args ...interface{}) interface{} {
			deReferencedFunc := *callback
			result := deReferencedFunc(args[0].(string), args[1].(containercontract.Container), args[2:]...)
			return result
		}
		c.resolvingCallbacks[key] = append(c.resolvingCallbacks[key], typedCallback)
	}
	return nil
}

func (c *Container) AfterResolving(abstract interface{}, callback *func(string, containercontract.Container, ...interface{}) error) error {
	var key string

	switch v := abstract.(type) {
	case string:
		key = c.getAlias(v)
	case func(...interface{}) error:
		typedCallback := func(args ...interface{}) interface{} {
			currentCallback := abstract.(func(...interface{}) interface{})
			result := currentCallback(args...)
			return result
		}
		c.globalAfterResolvingCallbacks = append(c.globalAfterResolvingCallbacks, typedCallback)
		return nil
	case func(...interface{}) interface{}:
		typedCallback := abstract.(func(...interface{}) interface{})
		c.globalAfterResolvingCallbacks = append(c.globalAfterResolvingCallbacks, typedCallback)
		return nil
	default:
		return fmt.Errorf("abstract must be a string or a function")
	}

	if callback != nil {
		typedCallback := func(args ...interface{}) interface{} {
			deReferencedFunc := *callback
			result := deReferencedFunc(args[0].(string), args[1].(containercontract.Container), args[2:]...)
			return result
		}
		c.afterResolvingCallbacks[key] = append(c.afterResolvingCallbacks[key], typedCallback)
	}
	return nil
}

func (c *Container) Instance(abstract string, instance interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeAbstractAlias(abstract)

	isBound := c.Bound(abstract)

	delete(c.aliases, abstract)

	c.instances[abstract] = instance

	if isBound {
		c.rebound(abstract)
	}

	return instance, nil
}

func (c *Container) ForgetInstance(abstract string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	instance := c.instances[abstract]
	if instance != nil {
		delete(c.instances, abstract)
	}
}

func (c *Container) IsInstance(abstract string) bool {
	_, ok := c.instances[abstract]
	return ok
}

func (c *Container) resolveKeys(keys interface{}) ([]string, error) {
	switch keySlice := keys.(type) {
	case []string:
		return keySlice, nil
	case string:
		return []string{keySlice}, nil
	case []interface{}:
		strKeys := make([]string, len(keySlice))
		for i, key := range keySlice {
			strKey, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported key type: %T", key)
			}
			strKeys[i] = strKey
		}
		return strKeys, nil
	default:
		return nil, fmt.Errorf("unsupported key type: %T", keys)
	}
}

func (c *Container) isShared(abstract string) bool {
	if _, ok := c.instances[abstract]; ok {
		return true
	}

	absBinding, ok := c.bindings[abstract]
	if ok {
		return absBinding.shared
	}

	return false
}

func (c *Container) isAlias(abstract string) bool {
	_, ok := c.aliases[abstract]
	return ok
}

func (c *Container) getAlias(abstract string) string {
	if alias, ok := c.aliases[abstract]; ok {
		return c.getAlias(alias)
	}

	return abstract
}

func (c *Container) Resolve(abstract string, raiseEvents bool, parameters ...interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if alias, ok := c.aliases[abstract]; ok {
		abstract = alias
	}

	if raiseEvents {
		c.fireBeforeResolvingCallbacks(abstract, parameters)
	}

	concrete, needsContextualBuild, err := c.getContextualConcrete(abstract, parameters)
	if err != nil {
		return nil, err
	}

	if instance, ok := c.instances[abstract]; ok && !needsContextualBuild {
		return instance, nil
	}

	var object interface{}
	if c.isBuildable(concrete, abstract) {
		object, err = c.build(concrete)
	} else {
		object, err = c.Resolve(concrete.(string), false, parameters)
	}

	if err != nil {
		return nil, err
	}

	for _, extender := range c.getExtenders(abstract) {
		object, _ = extender(object, c)
	}

	if c.isShared(abstract) && !needsContextualBuild {
		c.instances[abstract] = object
	}

	if raiseEvents {
		c.fireResolvingCallbacks(abstract, object)
	}

	c.resolved[abstract] = true

	return object, nil
}

func (c *Container) build(concrete interface{}) (interface{}, error) {
	concreteType := reflect.TypeOf(concrete)

	switch concreteType.Kind() {
	case reflect.Func:
		return c.buildFromFactory(concreteType, concrete)
	case reflect.Struct:
		return c.buildFromClass(concreteType)
	default:
		return nil, fmt.Errorf("unsupported concrete type: %T", concrete)
	}
}

func (c *Container) buildFromFactory(concreteType reflect.Type, concrete interface{}) (interface{}, error) {
	concreteValue := reflect.ValueOf(concrete)

	if concreteType.NumOut() != 1 {
		return nil, fmt.Errorf("factory function must return exactly one value")
	}

	if concreteType.NumIn() != 0 {
		return nil, fmt.Errorf("factory function must not have any input parameters")
	}

	instance := concreteValue.Call([]reflect.Value{})
	return instance[0].Interface(), nil
}

func (c *Container) buildFromClass(concreteType reflect.Type) (interface{}, error) {
	concreteValue := reflect.New(concreteType)

	if !concreteType.AssignableTo(concreteType) {
		return nil, containercontract.NewBindingResolutionError(
			fmt.Sprintf("Concrete type %s is not instantiable.", concreteType),
		)
	}

	return concreteValue.Interface(), nil
}

func (c *Container) dropStaleInstances(abstract string) {
	delete(c.instances, abstract)
	delete(c.aliases, abstract)
}

func (c *Container) getClosure(abstract string, concrete interface{}) func(container *Container, parameters ...interface{}) (interface{}, error) {
	return func(container *Container, parameters ...interface{}) (interface{}, error) {
		if abstract == concrete {
			return container.build(concrete)
		}

		return c.Resolve(concrete.(string), true, parameters...)
	}
}

func (c *Container) rebound(abstract string) {
	instance, err := c.Make(abstract)
	if err != nil {
		return // You can decide how to handle this error, e.g. log it or return it from the function
	}

	for _, callback := range c.getReboundCallbacks(abstract) {
		callback(c, instance)
	}
}

func (c *Container) getReboundCallbacks(abstract string) []func(container containercontract.Container, instance interface{}) {
	if callbacks, ok := c.reboundCallbacks[abstract]; ok {
		return callbacks
	}
	return []func(container containercontract.Container, instance interface{}){}
}

func (c *Container) fireBeforeResolvingCallbacks(abstract string, parameters ...interface{}) {
	c.fireBeforeCallbackArray(abstract, c.globalBeforeResolvingCallbacks, parameters)

	for typ, callbacks := range c.beforeResolvingCallbacks {
		if typ == abstract {
			c.fireBeforeCallbackArray(abstract, callbacks, parameters)
		}
	}
}

func (c *Container) fireBeforeCallbackArray(abstract string, callbacks []func(...interface{}) interface{}, parameters ...interface{}) {
	for _, callback := range callbacks {
		callback(abstract, c, parameters)
	}
}

func (c *Container) fireResolvingCallbacks(abstract string, object interface{}) {
	c.fireCallbackArray(object, c.globalResolvingCallbacks)

	for typ, callbacks := range c.resolvingCallbacks {
		if typ == abstract {
			c.fireCallbackArray(object, callbacks)
		}
	}
}

func (c *Container) fireCallbackArray(object interface{}, callbacks []func(...interface{}) interface{}) {
	for _, callback := range callbacks {
		callback(object, c)
	}
}

func (c *Container) getContextualConcrete(abstract string, parameters ...interface{}) (interface{}, bool, error) {
	if concrete, ok := c.contextualBindings[abstract]; ok {
		if !isEmpty(parameters) {
			return concrete, true, nil
		}
		return concrete, false, nil
	}
	return nil, false, fmt.Errorf("no contextual concrete found for %s", abstract)
}

func isEmpty(s []interface{}) bool {
	return len(s) == 0
}

func isFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

func (c *Container) isBuildable(concrete interface{}, abstract string) bool {
	if concrete.(string) == abstract {
		return true
	}
	return isFunc(concrete)
}

func (c *Container) getExtenders(abstract string) []func(...interface{}) (interface{}, error) {
	abstractAlias := c.getAlias(abstract)
	if extenders, ok := c.extenders[abstractAlias]; ok {
		return extenders
	}
	return []func(...interface{}) (interface{}, error){}
}

func (c *Container) forgetExtenders(abstract string) {
	abstractAlias := c.getAlias(abstract)
	delete(c.extenders, abstractAlias)
}

func (c *Container) removeAbstractAlias(searched string) {
	if _, exists := c.aliases[searched]; !exists {
		return
	}

	for abstract, aliases := range c.abstractAliases {
		for index, alias := range aliases {
			if alias == searched {
				c.abstractAliases[abstract] = append(c.abstractAliases[abstract][:index], c.abstractAliases[abstract][index+1:]...)
				break
			}
		}
	}
}
