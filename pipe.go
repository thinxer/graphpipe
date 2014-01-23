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

	next *uniqueue

	verbose bool
}

// Return the id of next tick.
func (p *GraphPipe) TickId() int {
	return p.tid
}

func (p *GraphPipe) startSources() {
	p.next = newUniqueue(len(p.nodes))
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
	// start sources if not already
	if !p.started {
		p.startSources()
		p.started = true
	}

	p.tid++
	if p.verbose {
		log.Printf("GraphPipe[%d] started.", p.tid)
	}

	// check if anything scheduled from last run.
	queue := p.next
	p.next = newUniqueue(len(p.nodes))
	// if not, check for signals.
	if queue.Len() == 0 {
		id, ok := <-p.control
		if !ok {
			if p.verbose {
				log.Printf("GraphPipe[%d] all sources closed.", p.tid)
			}
			return false
		}
		queue.Push(id)
	}

	// run the queue
	for queue.Len() > 0 {
		i := queue.Pop()
		uresult := p.nodes[i].Update(p.tid)
		if uresult&Update > 0 {
			for _, j := range p.children[i] {
				queue.Push(j)
			}
		}
		if uresult&More > 0 {
			p.next.Push(i)
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
