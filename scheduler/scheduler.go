package scheduler

import (
	"context"
	"time"
)

// Task represents a unit of work to be executed with specific priority and constraints.
type Task struct {
	ID       string
	Priority int
	ExecFunc func() error
	Timeout  time.Duration
}

// Scheduler manages the execution of Tasks using a pool of worker goroutines.
// It handles task submission, concurrent task processing, and graceful shutdown.
// tasks is a buffered channel used to queue Tasks for workers to process.
// workers specifies the number of worker goroutines for concurrent task processing.
// stopChan is a control channel to signal and handle graceful termination of workers.
type Scheduler struct {
	tasks    chan Task
	workers  int
	stopChan chan struct{}
}

// NewScheduler creates and starts a new Scheduler with the specified number of worker goroutines.
func NewScheduler(workers int) *Scheduler {
	s := &Scheduler{
		tasks:    make(chan Task, 100), //making room for 100 tasks
		workers:  workers,
		stopChan: make(chan struct{}),
	}
	s.Start()
	return s
}

func (s *Scheduler) Start() {
	for i := 0; i < s.workers; i++ {
		go func(workerID int) {
			for {
				select {
				case task := <-s.tasks:
					ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
					defer cancel()
					done := make(chan error, 1)
					go func() {
						done <- task.ExecFunc()
					}()
					select {
					case err := <-done:
						if err != nil {
							println("Worker", workerID, "failed task", task.ID, ":", err.Error())
						}
					case <-ctx.Done():
						println("Worker", workerID, "timed out task", task.ID)
					}
				case <-s.stopChan:
					return
				}
			}
		}(i)
	}
}

// Submit adds a Task to the Scheduler's task queue for processing by available workers.
func (s *Scheduler) Submit(task Task) {
	s.tasks <- task
}

// Stop gracefully shuts down the Scheduler by closing the stop channel, signaling all workers to terminate.
func (s *Scheduler) Stop() {
	close(s.stopChan) //Shut er down
}
