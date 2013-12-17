package graphpipe

import (
	"log"
	"reflect"
)

var nodeInterface = reflect.TypeOf((*Node)(nil)).Elem()

// the value of this registry will be reflected.
type registry map[string]interface{}

var defaultRegistry registry = registry(make(map[string]interface{}))

func (r registry) Regsiter(name string, newfunc interface{}) {
	t := reflect.TypeOf(newfunc)
	if t.Kind() != reflect.Func {
		panic("newfunc is not a func")
	}
	if t.NumIn() < 1 {
		panic("newfunc must have >=1 inputs")
	}
	if t.NumOut() != 1 {
		panic("newfunc must have exactly 1 output")
	}
	configT := t.In(0)
	if configT.Kind() != reflect.Ptr {
		panic("newfunc's first input must be a pointer (of config)")
	}
	returnT := t.Out(0)
	if returnT.Kind() != reflect.Ptr {
		panic("newfunc must return a pointer")
	}
	if !returnT.Implements(nodeInterface) {
		panic("newfunc must return a Node")
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

func (r registry) NewNode(name string, config interface{}, deps ...Node) Node {
	newfunc, ok := r[name]
	if !ok {
		log.Panicf("Node of %s not found", name)
	}

	ins := []reflect.Value{reflect.ValueOf(config)}
	for _, dep := range deps {
		ins = append(ins, reflect.ValueOf(dep))
	}
	outs := reflect.ValueOf(newfunc).Call(ins)
	return outs[0].Interface().(Node)
}

func Regsiter(name string, newfunc interface{}) {
	defaultRegistry.Regsiter(name, newfunc)
}

func NewConfig(name string) interface{} {
	return defaultRegistry.NewConfig(name)
}

func NewNode(name string, config interface{}, deps ...Node) Node {
	return defaultRegistry.NewNode(name, config, deps...)
}
