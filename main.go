package main

import (
	"fmt"
	"goConcurrencyTaskScheduler/scheduler"
	"time"
)

func main() {

	s := scheduler.NewScheduler(1)
	time.Sleep(50 * time.Millisecond)

	for _, t := range []struct {
		id       string
		priority int
	}{
		{"p10", 10},
		{"p1", 1},
		{"p20", 20},
		{"p5", 5},
	} {
		id := t.id
		priority := t.priority

		fmt.Println("SUBMIT", id, "priority", priority)
		s.Submit(scheduler.Task{
			ID:       id,
			Priority: priority,
			Timeout:  2 * time.Second,
			ExecFunc: func() error {
				fmt.Println("RUN", id)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		})
	}

	time.Sleep(3 * time.Second)
	s.Stop()
}
