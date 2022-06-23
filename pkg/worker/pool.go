package worker

import (
	"fmt"
	"sync"
	"time"
)

// Pool is an interface for worker pools
type Pool interface {
	Schedule(interface{ Execute() }) error
	Close()
	Wait()
}

const checkClosedInterval = time.Millisecond * 500

// FinitePool contains logic of goroutine reuse
type FinitePool struct {
	wg    *sync.WaitGroup
	size  int
	stop  bool
	tasks chan interface{ Execute() }
	done  chan struct{}
}

// New creates new worker pool with pre-defined amount of spawned goroutines and fixed queue size.
func New(nWorkers int, queueSize int) *FinitePool {
	pool := &FinitePool{
		wg:    &sync.WaitGroup{},
		tasks: make(chan interface{ Execute() }, queueSize),
		done:  make(chan struct{}),
	}

	for pool.size < nWorkers {
		pool.size++
		pool.wg.Add(1)

		go pool.worker()
	}

	return pool
}

// Schedule task to be executed
func (p *FinitePool) Schedule(task interface{ Execute() }) error {
	if p.stop {
		return fmt.Errorf("pool closed")
	}

	p.tasks <- task

	return nil
}

func (p *FinitePool) worker() {
	defer p.wg.Done()

	for {
		select {
		case task, ok := <-p.tasks:
			if !ok { // channel is closed + empty
				return
			}

			task.Execute()
		case <-p.done:
			return
		}
	}
}

// Close will close tasks queue
func (p *FinitePool) Close() {
	if p.stop {
		return
	}

	// stop scheduling new tasks
	p.stop = true

	// but close tasks channel only when
	// it does not contain any task
	go func() {
		ticker := time.NewTicker(checkClosedInterval)

		for {
			select {
			case <-ticker.C:

				if len(p.tasks) > 0 {
					continue
				}

				close(p.tasks)

				return
			}
		}
	}()
}

// Wait for all workers to finish
func (p *FinitePool) Wait() {
	p.wg.Wait()
}
