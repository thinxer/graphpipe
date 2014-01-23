package graphpipe_example

import (
	"github.com/thinxer/graphpipe"
)

func ExampleIntSampler() {
	var yamlData = `
# This is a sample YAML configuration.
---
  verbose: no
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
      input:
        - fib
      config:
        interval: 3
    - name: logger
      type: IntLogger
      input:
        - fib
      config:
        name: BeforeSampling
        silent: yes
    - name: logger2
      type: IntLogger
      input:
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
	// AfterSampling[1]: 1[1]
	// AfterSampling[4]: 3[4]
}

func ExampleIntDiffer() {
	var yamlData = `
# This is a sample YAML configuration.
---
  verbose: no
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
      input:
        - fib
      config:
        delay: 2
    - name: differ
      type: IntDiffer
      input:
        - fib
        - delayed
    - type: IntLogger
      input:
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
	// Diff[1]: 1[1]
	// Diff[2]: 1[2]
	// Diff[3]: 1[3]
	// Diff[4]: 2[4]
	// Diff[5]: 3[5]
	// Diff[6]: 5[6]
	// Diff[7]: 3[7]
	// Diff[8]: 0[8]
}
