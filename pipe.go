package graphpipe

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"log"
)

type YAMLConfig struct {
	Nodes []*struct {
		Name     string
		Type     string
		Source   bool
		Requires []string
		Config   interface{}
	}
}

func readConfig(bytes []byte) (config *YAMLConfig, err error) {
	config = &YAMLConfig{}
	// first pass to get node types
	err = goyaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	for _, node := range config.Nodes {
		// Is there a way not to remarshal the config?
		if node.Config != nil {
			remarshal, err := goyaml.Marshal(node.Config)
			if err != nil {
				return nil, err
			}

			node.Config = NewConfig(node.Type)
			err = goyaml.Unmarshal(remarshal, node.Config)
			if err != nil {
				return nil, err
			}
		} else {
			node.Config = NewConfig(node.Type)
		}
	}
	return
}

type GraphPipe struct {
	tid      int
	nodes    []Node
	source   []bool
	children [][]int
}

// Construct a graphpipe from a YAML config.
func GraphPipeFromYAML(yaml []byte) (*GraphPipe, error) {
	config, err := readConfig(yaml)
	if err != nil {
		return nil, err
	}
	for i, n := range config.Nodes {
		log.Printf("Node[%d]: %v, config: %+v", i, n, n.Config)
	}
	ncount := len(config.Nodes)
	pipe := &GraphPipe{
		nodes:    make([]Node, ncount),
		source:   make([]bool, ncount),
		children: make([][]int, ncount),
	}
	nodesMap := make(map[string]int)
	hasSource := false
	for i, nodeConfig := range config.Nodes {
		var deps []Node
		for _, depsName := range nodeConfig.Requires {
			depIndex := nodesMap[depsName]
			dep := pipe.nodes[depIndex]
			deps = append(deps, dep)
			pipe.children[depIndex] = append(pipe.children[depIndex], i)
		}
		node, err := NewNode(nodeConfig.Type, nodeConfig.Config, deps...)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, fmt.Errorf("Create node %s failed", nodeConfig.Name)
		}

		pipe.nodes[i] = node
		pipe.source[i] = nodeConfig.Source
		hasSource = hasSource || nodeConfig.Source
		nodesMap[nodeConfig.Name] = i
	}
	if !hasSource {
		return nil, fmt.Errorf("You must specify at least one source node!")
	}
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
	log.Printf("GraphPipe[%d] started.", p.tid)
	activated := make([]bool, len(p.nodes))
	for i, node := range p.nodes {
		if activated[i] || (p.source[i] && !p.nodes[i].Closed()) {
			updated := node.Update(p.tid)
			if updated {
				for _, j := range p.children[i] {
					activated[j] = true
				}
			}
			if node.Closed() {
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
	log.Printf("GraphPipe[%d] finished.", p.tid)
	p.tid++
	return closed < len(p.nodes)
}

// Run and empty the pipe.
func (p *GraphPipe) Run() {
	for p.RunOnce() {
		// Just loop till the end of the world.
	}
}
