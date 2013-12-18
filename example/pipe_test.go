package graphpipe_example

import (
	"github.com/thinxer/graphpipe"
)

func ExampleGraphpipe() {
	var yamlData = `
# This is a sample YAML configuration.
---
  nodes:
    - name: fib
      type: Fibonacci
      source: yes
      config:
        seed1: 0
        seed2: 1
        limit: 5
    - name: sampler
      type: IntSampler
      requires:
        - fib
      config:
        interval: 3
    - name: logger
      type: IntLogger
      requires:
        - fib
      config:
        name: BeforeSampling
        silent: yes
    - name: logger2
      type: IntLogger
      requires:
        - sampler
      config:
        name: AfterSampling
`
	pipe, err := graphpipe.GraphPipeFromYAML([]byte(yamlData))
	if err != nil {
		panic(err)
	}
	pipe.Run()

	// BeforeSampling[0]: 1
	// BeforeSampling[1]: 1
	// BeforeSampling[2]: 2
	// BeforeSampling[3]: 3
	// BeforeSampling[4]: 5

	// Output:
	// AfterSampling[0]: 1
	// AfterSampling[3]: 3
}
