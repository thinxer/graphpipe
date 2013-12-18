package graphpipe_example

import (
	"github.com/thinxer/graphpipe"
)

func ExampleIntSampler() {
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

	// Output:
	// AfterSampling[0]: 1[0]
	// AfterSampling[3]: 3[3]
}

func ExampleIntDiffer() {
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
        limit: 6
    - name: delayed
      type: IntDelayer
      requires:
        - fib
      config:
        delay: 2
    - name: differ
      type: IntDiffer
      requires:
        - fib
        - delayed
    - name: logger
      type: IntLogger
      requires:
        - fib
      config:
        name: Fib
        silent: yes
    - name: logger2
      type: IntLogger
      requires:
        - delayed
      config:
        name: Delayed
        silent: yes
    - name: logger3
      type: IntLogger
      requires:
        - differ
      config:
        name: Diff
`
	pipe, err := graphpipe.GraphPipeFromYAML([]byte(yamlData))
	if err != nil {
		panic(err)
	}
	pipe.Run()

	// Output:
	// Diff[0]: 1[0]
	// Diff[1]: 1[1]
	// Diff[2]: 1[2]
	// Diff[3]: 2[3]
	// Diff[4]: 3[4]
	// Diff[5]: 5[5]
}
