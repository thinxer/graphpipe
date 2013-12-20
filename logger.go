package graphpipe

import (
	"fmt"
	"reflect"
)

// A universal logger.
type Logger struct {
	name   string
	source AnySource
}

func (l *Logger) Update(mytid int) bool {
	valueMethod := reflect.ValueOf(l.source).MethodByName("Value")
	results := valueMethod.Call([]reflect.Value{})
	tid := results[0].Int()
	value := results[1].Interface()
	fmt.Printf("%s[%d]: %v[%d]\n", l.name, mytid, value, tid)
	return false
}

func (l *Logger) Value() int {
	return -1
}

func (l *Logger) Closed() bool {
	return l.source.Closed()
}

func NewLogger(config *struct{ Name string }, source AnySource) *Logger {
	valueMethod := reflect.ValueOf(source).MethodByName("Value")
	if valueMethod.Kind() != reflect.Func {
		panic("Cannot find Value method!")
	}
	if valueMethod.Type().NumIn() != 0 {
		panic("Value method must have 0 inputs!")
	}
	if valueMethod.Type().NumOut() != 2 {
		panic("Value method must have 2 outputs!")
	}
	if valueMethod.Type().Out(0).Kind() != reflect.Int {
		panic("Value method must return (int, _)!")
	}
	return &Logger{name: config.Name, source: source}
}

func init() {
	Regsiter("Logger", NewLogger)
}
