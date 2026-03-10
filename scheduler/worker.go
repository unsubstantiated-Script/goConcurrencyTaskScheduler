package scheduler

import (
	"context"
	"fmt"
	"time"
)

// startWorker runs a worker loop that continuously retrieves tasks from the task channel and executes them.
func (s *Scheduler) startWorker(id int) {
	for {
		select {
		case task := <-s.taskChan:
			runTaskWithTimeout(id, task)
		case <-s.stopChan:
			return
		}
	}
}

// runTaskWithTimeout executes the given task with a timeout.
func runTaskWithTimeout(workerID int, task Task) {
	timeout := task.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- task.ExecFunc()
	}()

	select {
	case err := <-done:
		if err != nil {
			fmt.Println("worker", workerID, "task failed:", task.ID, err)
		}
	case <-ctx.Done():
		fmt.Println("worker", workerID, "task timed out:", task.ID)
	}
}
