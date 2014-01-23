package graphpipe

import (
	"sync"

	"log"
)

type GraphPipe struct {
	tid      int
	nodes    []Node
	children [][]int

	control chan int
	started bool
	wg      sync.WaitGroup

	verbose bool
}

// Return the id of next tick.
func (p *GraphPipe) TickId() int {
	return p.tid
}

func (p *GraphPipe) startSources() {
	p.control = make(chan int, 128)
	for id, node := range p.nodes {
		if source, ok := node.(SourceNode); ok {
			if p.verbose {
				log.Printf("Starting source %d", id)
			}
			p.wg.Add(1)
			ch := make(chan bool, 128)
			// feed
			go source.Start(ch)
			// trans feed
			go func(id int) {
				for _ = range ch {
					p.control <- id
				}
				p.wg.Done()
			}(id)
		}
	}
	go func() {
		p.wg.Wait()
		close(p.control)
	}()
}

// RunOnce will select an active source and start the pipe.
// tid may be increased one or more,
// depending on whether there will be HasMore updates.
func (p *GraphPipe) RunOnce() bool {
	if !p.started {
		p.startSources()
		p.started = true
	}

	queue := make([]int, 0, len(p.nodes))
	queued := make([]bool, len(p.nodes))
	enqueue := func(i int) {
		if queued[i] {
			return
		}
		queue = append(queue, i)
		queued[i] = true
	}
	dequeue := func() int {
		i := queue[0]
		queue = queue[1:]
		queued[i] = false
		return i
	}

	seeds := make([]bool, len(p.nodes))
	activatedSource, ok := <-p.control
	if !ok {
		if p.verbose {
			log.Printf("GraphPipe[%d] all sources closed.", p.tid)
		}
		return false
	}
	seeds[activatedSource] = true

	more := true
	for more {
		more = false
		p.tid++
		if p.verbose {
			log.Printf("GraphPipe[%d] started.", p.tid)
		}

		for i, s := range seeds {
			if s {
				enqueue(i)
				seeds[i] = false
			}
		}
		for len(queue) > 0 {
			i := dequeue()
			node := p.nodes[i]
			updated := node.Update(p.tid)
			if updated != Skip || node.Closed() {
				for _, j := range p.children[i] {
					enqueue(j)
				}
				if updated == HasMore {
					seeds[i] = true
					more = true
				}
			}

			if p.verbose && node.Closed() {
				log.Printf("GraphPipe[%d] Node[%d] Closed", p.tid, i)
			}
		}
	}
	if p.verbose {
		log.Printf("GraphPipe[%d] finished.", p.tid)
	}
	return true
}

// Run and empty the pipe.
func (p *GraphPipe) Run() {
	for p.RunOnce() {
		// Just loop till the end of the world.
	}
}
