package graphpipe

import (
	"fmt"
	"io/ioutil"
	"log"

	"launchpad.net/goyaml"
)

type NodeConfig struct {
	Name   string
	Type   string
	Source bool
	Inject []string
	Input  []string
	Config interface{}
}

type Config struct {
	Verbose         bool
	Services, Nodes []NodeConfig
}

func readConfig(bytes []byte) (config *Config, err error) {
	config = &Config{}
	// first pass to get node types
	err = goyaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	// second pass to get node configs
	for _, nodes := range [][]NodeConfig{config.Services, config.Nodes} {
		for i := range nodes {
			// Is there a way not to remarshal the config?
			if nodes[i].Config != nil {
				remarshal, err := goyaml.Marshal(nodes[i].Config)
				if err != nil {
					return nil, err
				}

				nodes[i].Config = NewConfig(nodes[i].Type)
				err = goyaml.Unmarshal(remarshal, nodes[i].Config)
				if err != nil {
					return nil, err
				}
			} else {
				nodes[i].Config = NewConfig(nodes[i].Type)
			}
		}
	}
	return
}

type GraphPipe struct {
	tid      int
	nodes    []Node
	source   []bool
	children [][]int

	verbose bool
}

// Construct a graphpipe from a YAML config.
func GraphPipeFromYAML(yaml []byte) (*GraphPipe, error) {
	config, err := readConfig(yaml)
	if err != nil {
		return nil, err
	}
	ncount := len(config.Nodes)
	pipe := &GraphPipe{
		nodes:    make([]Node, ncount),
		source:   make([]bool, ncount),
		children: make([][]int, ncount),
	}

	// services map
	servicesMap := make(map[string]int)
	services := make([]interface{}, len(config.Services))
	newNode := func(nodeConfig NodeConfig) (interface{}, error) {
		if config.Verbose {
			log.Printf("Creating: %+v\n", nodeConfig)
		}
		var injects []interface{}
		for _, serviceName := range nodeConfig.Inject {
			serviceIndex := servicesMap[serviceName]
			service := services[serviceIndex]
			injects = append(injects, service)
		}
		something, err := NewNode(nodeConfig.Type, nodeConfig.Config, injects...)

		if something == nil && err == nil {
			err = fmt.Errorf("Creation of %s failed", nodeConfig.Name)
		}
		return something, err
	}

	// setup Services
	for i, serviceConfig := range config.Services {
		service, err := newNode(serviceConfig)
		if err != nil {
			return nil, err
		}
		services[i] = service
		servicesMap[serviceConfig.Name] = i
	}

	// setup Nodes
	nodesMap := make(map[string]int)
	hasSource := false
	for i, nodeConfig := range config.Nodes {
		nodeV, err := newNode(nodeConfig)
		if err != nil {
			return nil, err
		}
		node := nodeV.(Node)
		if len(nodeConfig.Input) > 0 {
			var sources []Node
			for _, nodeName := range nodeConfig.Input {
				depIndex := nodesMap[nodeName]
				dep := pipe.nodes[depIndex]
				sources = append(sources, dep)
				pipe.children[depIndex] = append(pipe.children[depIndex], i)
			}
			if err := SetInput(node, sources...); err != nil {
				return nil, err
			}
		}

		pipe.nodes[i] = node
		pipe.source[i] = nodeConfig.Source
		hasSource = hasSource || nodeConfig.Source
		nodesMap[nodeConfig.Name] = i
	}

	if !hasSource {
		return nil, fmt.Errorf("You must specify at least one source node!")
	}
	pipe.verbose = config.Verbose
	return pipe, nil
}

func GraphPipeFromYAMLFile(filename string) (*GraphPipe, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return GraphPipeFromYAML(bytes)
}

// Return the id of next tick.
func (p *GraphPipe) TickId() int {
	return p.tid
}

// Run once and increase tid by one.
func (p *GraphPipe) RunOnce() bool {
	closed := 0
	if p.verbose {
		log.Printf("GraphPipe[%d] started.", p.tid)
	}
	activated := make([]bool, len(p.nodes))
	for i, node := range p.nodes {
		if activated[i] || (p.source[i] && !p.nodes[i].Closed()) {
			updated := node.Update(p.tid)
			if updated {
				for _, j := range p.children[i] {
					activated[j] = true
				}
			}
			if node.Closed() && p.verbose {
				log.Printf("GraphPipe[%d] Node[%d] Closed", p.tid, i)
			}
		} else if p.nodes[i].Closed() {
			for _, j := range p.children[i] {
				if !p.nodes[j].Closed() {
					activated[j] = true
				}
			}
		}
		if node.Closed() {
			closed++
		}
	}
	if p.verbose {
		log.Printf("GraphPipe[%d] finished.", p.tid)
	}
	p.tid++
	return closed < len(p.nodes)
}

// Run and empty the pipe.
func (p *GraphPipe) Run() {
	for p.RunOnce() {
		// Just loop till the end of the world.
	}
}
