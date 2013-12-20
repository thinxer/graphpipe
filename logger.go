package graphpipe

import (
	"fmt"
	"reflect"
)

// A universal logger.
type Logger struct {
	name    string
	sources []AnySource
}

func (l *Logger) Update(mytid int) bool {
	fmt.Printf("%s[%d]:", l.name, mytid)
	for _, source := range l.sources {
		valueMethod := reflect.ValueOf(source).MethodByName("Value")
		results := valueMethod.Call([]reflect.Value{})
		tid := results[0].Int()
		value := results[1].Interface()
		fmt.Printf("\t%v[%d]", value, tid)
	}
	fmt.Println()
	return false
}

func (l *Logger) Value() int {
	return -1
}

func (l *Logger) Closed() bool {
	for _, s := range l.sources {
		if !s.Closed() {
			return false
		}
	}
	return true
}

func NewLogger(config *struct{ Name string }, sources ...AnySource) *Logger {
	for i, source := range sources {
		valueMethod := reflect.ValueOf(source).MethodByName("Value")
		if valueMethod.Kind() != reflect.Func {
			fmt.Errorf("%d source: Cannot find Value method!", i)
		}
		if valueMethod.Type().NumIn() != 0 {
			fmt.Errorf("%d source: Value method must have 0 inputs!", i)
		}
		if valueMethod.Type().NumOut() != 2 {
			fmt.Errorf("%d source: Value method must have 2 outputs!", i)
		}
		if valueMethod.Type().Out(0).Kind() != reflect.Int {
			fmt.Errorf("%d source: Value method must return (int, _)!", i)
		}
	}
	return &Logger{name: config.Name, sources: sources}
}

func init() {
	Regsiter("Logger", NewLogger)
}
