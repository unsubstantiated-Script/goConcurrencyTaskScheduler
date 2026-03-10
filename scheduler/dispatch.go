package scheduler

import (
	"container/heap"
	"time"
)

func (s *Scheduler) dispatch() {
	for {
		select {
		case <-s.stopChan:
			return
		default:
		}

		s.mu.Lock()
		if s.tasks.Len() == 0 {
			s.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		task := heap.Pop(s.tasks).(*Task)
		s.mu.Unlock()

		select {
		case s.taskChan <- *task:
		case <-s.stopChan:
			return
		}
	}
}
