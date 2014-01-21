package graphpipe

import "log"

type GraphPipe struct {
	tid      int
	nodes    []Node
	source   []bool
	children [][]int

	verbose bool
}

// Return the id of next tick.
func (p *GraphPipe) TickId() int {
	return p.tid
}

// Run once and increase tid by one.
func (p *GraphPipe) RunOnce() bool {
	if p.verbose {
		log.Printf("GraphPipe[%d] started.", p.tid)
	}

	seeds := make([]bool, len(p.nodes))
	copy(seeds, p.source)

	queue := []int{}
	activated := make([]bool, len(p.nodes))
	enqueue := func(i int) {
		if activated[i] {
			return
		}
		queue = append(queue, i)
		activated[i] = true
	}

	closed := true
	for {
		for i, s := range seeds {
			if s {
				enqueue(i)
				seeds[i] = false
			}
		}
		if len(queue) == 0 {
			break
		}
		for len(queue) > 0 {
			i := queue[0]
			queue = queue[1:]
			activated[i] = false
			node := p.nodes[i]
			if node.Closed() {
				continue
			}

			updated := node.Update(p.tid)
			if updated != Skip || (p.source[i] && node.Closed()) {
				for _, j := range p.children[i] {
					if !activated[j] {
						enqueue(j)
					}
				}
				if updated == HasMore {
					seeds[i] = true
				}
			}

			if !node.Closed() {
				closed = false
			}

			if p.verbose && node.Closed() {
				log.Printf("GraphPipe[%d] Node[%d] Closed", p.tid, i)
			}
		}
	}
	if p.verbose {
		log.Printf("GraphPipe[%d] finished.", p.tid)
	}

	p.tid++
	return !closed
}

// Run and empty the pipe.
func (p *GraphPipe) Run() {
	for p.RunOnce() {
		// Just loop till the end of the world.
	}
}
