package foundation

import (
	"errors"
	"fmt"
	"github.com/goravel/framework/contracts/config"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"reflect"
	"strings"
	"sync"
)

type BaseManager struct {
	app            foundationcontract.Application
	config         config.Config
	drivers        map[string]interface{}
	customCreators map[string]func() interface{}
	mu             sync.RWMutex
}

func NewBaseManager(app *Application) *BaseManager {
	configManager, err := app.Make("config")
	if err != nil {
		panic(err)
	}
	return &BaseManager{
		app:            app,
		config:         configManager.(config.Config),
		drivers:        make(map[string]interface{}),
		customCreators: make(map[string]func() interface{}),
	}
}

func (bm *BaseManager) GetDefaultDriver() string {
	// Implement logic to get the default driver name.
	return ""
}

func (bm *BaseManager) Driver(driver string) (interface{}, error) {
	if driver == "" {
		driver = bm.GetDefaultDriver()
	}
	bm.mu.RLock()
	instance, ok := bm.drivers[driver]
	bm.mu.RUnlock()
	if ok {
		return instance, nil
	}

	if creator, ok := bm.customCreators[driver]; ok {
		instance = creator()
	} else {
		methodName := "Create" + strings.Title(driver) + "Driver"
		method := reflect.ValueOf(bm).MethodByName(methodName)
		if !method.IsValid() {
			return nil, errors.New("driver not supported")
		}
		instance = method.Call(nil)[0].Interface()
	}

	bm.mu.Lock()
	bm.drivers[driver] = instance
	bm.mu.Unlock()

	return instance, nil
}

func (bm *BaseManager) Extend(driver string, creator func() interface{}) {
	bm.customCreators[driver] = creator
}

func (bm *BaseManager) GetDrivers() map[string]interface{} {
	return bm.drivers
}

func (bm *BaseManager) GetContainer() foundationcontract.Application {
	return bm.app
}

func (bm *BaseManager) SetContainer(container foundationcontract.Application) {
	bm.app = container
}

func (bm *BaseManager) ForgetDrivers() {
	bm.mu.Lock()
	bm.drivers = make(map[string]interface{})
	bm.mu.Unlock()
}

func (bm *BaseManager) Call(methodName string, args ...interface{}) (interface{}, error) {
	driver, err := bm.Driver("")
	if err != nil {
		return nil, err
	}

	methodValue := reflect.ValueOf(driver).MethodByName(methodName)
	if !methodValue.IsValid() {
		return nil, fmt.Errorf("method %s does not exist on driver %s", methodName, driver)
	}

	var input []reflect.Value
	for _, arg := range args {
		input = append(input, reflect.ValueOf(arg))
	}

	output := methodValue.Call(input)
	if len(output) == 0 {
		return nil, nil
	}

	return output[0].Interface(), nil
}
