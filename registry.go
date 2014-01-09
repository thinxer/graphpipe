package graphpipe

import (
	"fmt"
	"log"
	"reflect"
)

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

var registry = make(map[string]interface{})

// Register a newNode function to the default registry.
// Please do not export the newNode function.
// This function will examine the newNode func for necessary methods,
// and get the config type from the first argument.
func Regsiter(name string, newfunc interface{}) {
	t := reflect.TypeOf(newfunc)
	if t.Kind() != reflect.Func {
		panic(name + "'s newfunc is not a func")
	}
	if t.NumIn() < 1 {
		panic(name + "'s newfunc must have >=1 inputs")
	}
	if t.NumOut() != 2 {
		panic(name + "'s newfunc must return (_, error)")
	}
	configT := t.In(0)
	if configT.Kind() != reflect.Ptr {
		panic(name + "'s newfunc's first input must be a pointer (of config)")
	}
	if !t.Out(1).Implements(errorInterface) {
		panic(name + "'s newfunc must return an error as the second output")
	}
	registry[name] = newfunc
}

// Create a new config by node name from the default registry.
func NewConfig(name string) interface{} {
	newfunc, ok := registry[name]
	if !ok {
		log.Panicf("Type of %s not found", name)
	}
	configType := reflect.TypeOf(newfunc).In(0).Elem()
	return reflect.New(configType).Interface()
}

// Create a new node by name, config and dependencies from the default registry.
func NewNode(name string, config interface{}, services ...interface{}) (interface{}, error) {
	newfunc, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("Node of %s not found", name)
	}
	for _, service := range services {
		if service == nil {
			return nil, fmt.Errorf("Nil service detected")
		}
	}

	ins := []reflect.Value{reflect.ValueOf(config)}
	for _, service := range services {
		ins = append(ins, reflect.ValueOf(service))
	}
	outs := reflect.ValueOf(newfunc).Call(ins)
	node := outs[0].Interface()
	err := outs[1].Interface()
	if err != nil {
		return nil, err.(error)
	} else {
		return node, nil
	}
}

// List registered types.
func List() (ret []string) {
	for k := range registry {
		ret = append(ret, k)
	}
	return
}

// Set inputs to the node.
func SetInput(node Node, sources ...Node) error {
	method := reflect.ValueOf(node).MethodByName("SetInput")
	if !method.IsValid() {
		return fmt.Errorf("Node have no SetSources method!")
	}
	ins := []reflect.Value{}
	for _, source := range sources {
		ins = append(ins, reflect.ValueOf(source))
	}
	method.Call(ins)
	return nil
}
