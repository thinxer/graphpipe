package graphpipe

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

type syncWriter interface {
	Sync() error
	io.Writer
}

// A universal logger.
type Logger struct {
	name    string
	sources []AnySource
	output  syncWriter
}

type LoggerConfig struct {
	Name   string
	Output string
}

func (l *Logger) Update(mytid int) bool {
	fmt.Fprintf(l.output, "%v\t%s[%d]:", time.Now(), l.name, mytid)
	for _, source := range l.sources {
		valueMethod := reflect.ValueOf(source).MethodByName("Value")
		results := valueMethod.Call([]reflect.Value{})
		tid := results[0].Int()
		value := results[1].Interface()
		fmt.Fprintf(l.output, "\t%v[%d]", value, tid)
	}
	fmt.Fprintln(l.output)
	l.output.Sync()
	return false
}

func (l *Logger) Closed() bool {
	for _, s := range l.sources {
		if !s.Closed() {
			return false
		}
	}
	return true
}

func NewLogger(config *LoggerConfig, sources ...AnySource) *Logger {
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
	var output syncWriter
	switch config.Output {
	case "", "-":
		output = os.Stdout
	case "--":
		output = os.Stderr
	default:
		var err error
		output, err = os.OpenFile(config.Output, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}
	fmt.Fprintln(output, "--------", time.Now())
	output.Sync()
	return &Logger{name: config.Name, sources: sources, output: output}
}

func init() {
	Regsiter("Logger", NewLogger)
}
