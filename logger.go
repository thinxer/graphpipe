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
// It is not recommended to use this Node in high performance programs,
// as reflection is extensively used.
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
	if l.Closed() {
		return false
	}
	fmt.Fprintf(l.output, "%v [%d]%s:", time.Now().Format("0102 15:04:05"), mytid, l.name)
	for _, source := range l.sources {
		valueMethod := reflect.ValueOf(source).MethodByName("Value")
		results := valueMethod.Call([]reflect.Value{})
		fmt.Fprintf(l.output, "\t")
		for _, r := range results {
			fmt.Fprintf(l.output, "[%v]", r.Interface())
		}
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

func newLogger(config *LoggerConfig, sources ...AnySource) (*Logger, error) {
	for i, source := range sources {
		valueMethod := reflect.ValueOf(source).MethodByName("Value")
		if valueMethod.Kind() != reflect.Func {
			fmt.Errorf("%d source: Cannot find Value method!", i)
		}
		if valueMethod.Type().NumIn() != 0 {
			fmt.Errorf("%d source: Value method must have 0 inputs!", i)
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
			return nil, err
		}
	}
	fmt.Fprintln(output, "--------", time.Now())
	output.Sync()
	return &Logger{name: config.Name, sources: sources, output: output}, nil
}

func init() {
	Regsiter("Logger", newLogger)
}
