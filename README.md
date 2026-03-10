# goConcurrencyTaskScheduler

A lightweight, concurrent task scheduler written in Go. Tasks are submitted with a priority level and executed by a configurable worker pool in priority order, with per-task timeout support.

---

## Features

- **Priority-based scheduling** — Tasks are queued in a max-heap priority queue and dispatched highest-priority first.
- **Configurable worker pool** — Spin up as many concurrent workers as needed at scheduler creation time.
- **Per-task timeouts** — Each task carries its own `time.Duration` timeout; workers cancel execution automatically if the deadline is exceeded.
- **Thread-safe submission** — Tasks can be submitted from multiple goroutines simultaneously without data races.
- **Graceful shutdown** — `Stop()` signals all workers and the dispatcher to exit cleanly.

---

## Project Structure

```
goConcurrencyTaskScheduler/
├── main.go                      # Example entry point
└── scheduler/
    ├── task.go                  # Task struct definition
    ├── priority_queue.go        # heap.Interface implementation (max-heap by Priority)
    ├── scheduler.go             # Scheduler struct, NewScheduler, Submit, Stop
    ├── dispatch.go              # Dispatch loop — pops tasks and sends to workers
    └── worker.go                # Worker loop and timeout execution logic
```

---

## Requirements

- Go **1.24** or later

---

## Getting Started

### Clone & Run

```zsh
git clone <repo-url>
cd goConcurrencyTaskScheduler
go run main.go
```

### Expected Output

```
SUBMIT p10 priority 10
SUBMIT p1 priority 1
SUBMIT p20 priority 20
SUBMIT p5 priority 5
RUN p20
RUN p10
RUN p5
RUN p1
```

Tasks are executed in descending priority order regardless of submission order.

---

## Usage

### 1. Create a Scheduler

```go
s := scheduler.NewScheduler(4) // 4 concurrent workers
```

### 2. Define and Submit Tasks

```go
s.Submit(scheduler.Task{
    ID:       "my-task",
    Priority: 15,
    Timeout:  2 * time.Second,
    ExecFunc: func() error {
        fmt.Println("doing work")
        return nil
    },
})
```

### 3. Stop the Scheduler

```go
s.Stop()
```

---

## API Reference

### `scheduler.Task`

| Field      | Type             | Description                                           |
|------------|------------------|-------------------------------------------------------|
| `ID`       | `string`         | Unique identifier for the task                        |
| `Priority` | `int`            | Higher value = higher priority                        |
| `ExecFunc` | `func() error`   | The function to execute                               |
| `Timeout`  | `time.Duration`  | Max execution time (defaults to 5s if ≤ 0)           |

### `scheduler.NewScheduler(workers int) *Scheduler`

Creates and starts a new Scheduler with the given number of worker goroutines and an internal dispatch loop.

### `(*Scheduler).Submit(task Task)`

Pushes a task onto the priority queue in a thread-safe manner.

### `(*Scheduler).Stop()`

Closes the stop channel, causing all workers and the dispatcher to return gracefully.

---

## How It Works

1. **Submit** pushes a `Task` pointer onto a `container/heap` max-heap keyed by `Priority`.
2. The **dispatch loop** polls the heap every 100 ms (or immediately when tasks are present), pops the highest-priority task, and sends it to the shared `taskChan`.
3. An available **worker** picks up the task from `taskChan` and calls `runTaskWithTimeout`, which races the task's `ExecFunc` against a `context.WithTimeout` deadline.
4. Calling **Stop** closes `stopChan`; all select statements listening on it exit, terminating the dispatch loop and every worker goroutine.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

---

## License

This project is unlicensed. Add a `LICENSE` file to define distribution terms.

