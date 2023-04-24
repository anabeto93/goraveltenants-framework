package container

import (
	"reflect"
)

// arrayWrap takes an interface and returns a slice of interfaces containing
// the input value if it's not nil and not a slice.
func arrayWrap(value interface{}) []interface{} {
	if value == nil {
		return []interface{}{}
	}

	valueSlice, ok := value.([]interface{})
	if ok {
		return valueSlice
	}

	return []interface{}{value}
}

// unwrapIfClosure checks if the value is a function and calls it with the provided arguments.
func unwrapIfClosure(value interface{}, args ...interface{}) (interface{}, error) {
	fnValue := reflect.ValueOf(value)
	if fnValue.Kind() != reflect.Func {
		return value, nil
	}

	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}

	result := fnValue.Call(argValues)
	if len(result) == 0 {
		return nil, nil
	}

	return result[0].Interface(), nil
}
