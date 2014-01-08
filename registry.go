package graphpipe

import (
	"fmt"
	"log"
	"reflect"
)

var (
	nodeInterface  = reflect.TypeOf((*Node)(nil)).Elem()
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

type registry map[string]interface{}

var defaultRegistry registry = registry(make(map[string]interface{}))

func (r registry) Regsiter(name string, newfunc interface{}) {
	t := reflect.TypeOf(newfunc)
	if t.Kind() != reflect.Func {
		panic(name + "'s newfunc is not a func")
	}
	if t.NumIn() < 1 {
		panic(name + "'s newfunc must have >=1 inputs")
	}
	if t.NumOut() != 2 {
		panic(name + "'s newfunc must return (Node, error)")
	}
	configT := t.In(0)
	if configT.Kind() != reflect.Ptr {
		panic(name + "'s newfunc's first input must be a pointer (of config)")
	}
	returnT := t.Out(0)
	if returnT.Kind() != reflect.Ptr {
		panic(name + "'s newfunc must return a pointer")
	}
	if !returnT.Implements(nodeInterface) {
		panic(name + "'s newfunc must return a Node as first output")
	}
	if !t.Out(1).Implements(errorInterface) {
		panic(name + "'s newfunc must return an error as second output")
	}
	r[name] = newfunc
}

func (r registry) NewConfig(name string) interface{} {
	newfunc, ok := r[name]
	if !ok {
		log.Panicf("Node of %s not found", name)
	}
	configType := reflect.TypeOf(newfunc).In(0).Elem()
	return reflect.New(configType).Interface()
}

func (r registry) NewNode(name string, config interface{}, deps ...Node) (Node, error) {
	newfunc, ok := r[name]
	if !ok {
		return nil, fmt.Errorf("Node of %s not found", name)
	}
	for _, dep := range deps {
		if dep == nil {
			return nil, fmt.Errorf("Nil dependency detected")
		}
	}

	ins := []reflect.Value{reflect.ValueOf(config)}
	for _, dep := range deps {
		ins = append(ins, reflect.ValueOf(dep))
	}
	outs := reflect.ValueOf(newfunc).Call(ins)
	node := outs[0].Interface()
	err := outs[1].Interface()
	if err != nil {
		return nil, err.(error)
	} else {
		return node.(Node), nil
	}
}

func (r registry) List() (ret []string) {
	for k := range r {
		ret = append(ret, k)
	}
	return
}

// Register a newNode function to the default registry.
// Please do not export the newNode function.
// This function will examine the newNode func for necessary methods,
// and get the config type from the first argument.
func Regsiter(name string, newfunc interface{}) {
	defaultRegistry.Regsiter(name, newfunc)
}

// List registered types.
func List() []string {
	return defaultRegistry.List()
}

// Create a new config by node name from the default registry.
func NewConfig(name string) interface{} {
	return defaultRegistry.NewConfig(name)
}

// Create a new node by name, config and dependencies from the default registry.
func NewNode(name string, config interface{}, deps ...Node) (Node, error) {
	return defaultRegistry.NewNode(name, config, deps...)
}
