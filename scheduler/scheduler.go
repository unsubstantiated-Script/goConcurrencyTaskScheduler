package scheduler

import (
	"container/heap"
	"sync"
)

// Scheduler manages a pool of workers and a priority queue for scheduling tasks.
type Scheduler struct {
	tasks    *PriorityQueue
	taskChan chan Task
	workers  int
	stopChan chan struct{}
	mu       sync.Mutex
}

// NewScheduler creates a new Scheduler instance with the specified number of workers.
func NewScheduler(workers int) *Scheduler {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	s := &Scheduler{
		tasks:    &pq,
		taskChan: make(chan Task, 100),
		workers:  workers,
		stopChan: make(chan struct{}),
	}

	go s.dispatch()

	for i := 0; i < workers; i++ {
		go s.startWorker(i)
	}

	return s
}

// Submit adds a new task to the scheduler's priority queue in a thread-safe manner.
func (s *Scheduler) Submit(task Task) {
	s.mu.Lock()
	heap.Push(s.tasks, &task)
	s.mu.Unlock()
}

// Stop gracefully shuts down the scheduler by closing the stopChan, signaling workers and dispatch loop to terminate.
func (s *Scheduler) Stop() {
	close(s.stopChan)
}
