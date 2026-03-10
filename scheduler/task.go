package scheduler

import "time"

// Task represents a task to be executed by a worker.
type Task struct {
	ID       string
	Priority int
	ExecFunc func() error
	Timeout  time.Duration
}
