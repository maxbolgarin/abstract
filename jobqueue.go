package abstract

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maxbolgarin/lang"
)

// JobQueue manages a pool of workers that execute context-aware tasks without return values.
// It's optimized for fire-and-forget operations with task tracking and wait capabilities.
type JobQueue struct {
	workers  int
	tasks    chan func(ctx context.Context)
	stopChan chan struct{}
	wg       sync.WaitGroup

	logger lang.Logger

	isQueueStarted    atomic.Bool
	onFlyRunningTasks atomic.Int64
	tasksInQueue      atomic.Int64
	finishedTasks     atomic.Int64
	totalTasks        atomic.Int64
}

// NewJobQueue creates a new context-aware job queue with the specified number of workers and task queue capacity.
func NewJobQueue(workers int, queueCapacity int, logger ...lang.Logger) *JobQueue {
	if workers <= 0 {
		workers = 1
	}
	if queueCapacity <= 0 {
		queueCapacity = workers * 100
	}

	return &JobQueue{
		workers:  workers,
		tasks:    make(chan func(ctx context.Context), queueCapacity),
		stopChan: make(chan struct{}),
		logger:   lang.First(logger),
	}
}

// Start launches the worker goroutines.
func (q *JobQueue) Start(ctx context.Context) {
	if !q.isQueueStarted.CompareAndSwap(false, true) {
		return
	}
	q.wg.Add(q.workers)
	for range q.workers {
		lang.Go(q.logger, func() {
			q.worker(ctx)
		})
	}
}

// Shutdown signals all workers to stop after completing their current tasks.
// It waits for all in-flight and queued tasks to complete or until the context is done.
func (q *JobQueue) Shutdown(ctx context.Context) error {
	if !q.isQueueStarted.CompareAndSwap(true, false) {
		return nil
	}
	close(q.stopChan)
	close(q.tasks)

	// Wait for all workers to finish
	done := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// StopNoWait signals all workers to stop after completing their current tasks.
// It does not wait for them to complete.
func (q *JobQueue) StopNoWait() {
	if !q.isQueueStarted.CompareAndSwap(true, false) {
		return
	}
	close(q.stopChan)
	close(q.tasks)
}

// Submit adds a task to the queue and returns true if the task was accepted.
// Returns false if the queue is stopped or the context is done.
func (q *JobQueue) Submit(ctx context.Context, task func(ctx context.Context)) bool {
	if task == nil {
		return false
	}
	if !q.isQueueStarted.Load() {
		return false
	}

	select {
	case q.tasks <- task:
		q.totalTasks.Add(1)
		q.tasksInQueue.Add(1)
		return true

	case <-q.stopChan:
		return false

	case <-ctx.Done():
		return false
	}
}

// Wait blocks until all submitted tasks have been completed or the context is done.
// Returns nil if all tasks completed successfully, or context error if cancelled.
func (q *JobQueue) Wait(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		// Check if all tasks are done
		if q.tasksInQueue.Load() == 0 && q.onFlyRunningTasks.Load() == 0 {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Continue waiting
		}
	}
}

// worker is the goroutine that processes tasks.
func (q *JobQueue) worker(ctx context.Context) {
	defer q.wg.Done()

	for {
		select {
		case task, ok := <-q.tasks:
			if !ok {
				// Channel closed, drain remaining tasks
				return
			}
			q.tasksInQueue.Add(-1)

			q.onFlyRunningTasks.Add(1)
			task(ctx)
			q.onFlyRunningTasks.Add(-1)

			q.finishedTasks.Add(1)

		case <-ctx.Done():
			return
		case <-q.stopChan:
			// Stop signal received, but continue processing pending tasks
			// The tasks channel will be closed, causing the worker to exit after draining
		}
	}
}

// TasksInQueue returns the number of tasks in the queue waiting to be executed.
func (q *JobQueue) TasksInQueue() int {
	return int(q.tasksInQueue.Load())
}

// OnFlyRunningTasks returns the number of currently executing tasks.
func (q *JobQueue) OnFlyRunningTasks() int {
	return int(q.onFlyRunningTasks.Load())
}

// FinishedTasks returns the number of completed tasks.
func (q *JobQueue) FinishedTasks() int {
	return int(q.finishedTasks.Load())
}

// TotalTasks returns the total number of tasks submitted to the queue.
func (q *JobQueue) TotalTasks() int {
	return int(q.totalTasks.Load())
}

// IsQueueStarted returns true if the job queue has been started.
func (q *JobQueue) IsQueueStarted() bool {
	return q.isQueueStarted.Load()
}

// PendingTasks returns the total number of tasks that are either queued or running.
func (q *JobQueue) PendingTasks() int {
	return int(q.tasksInQueue.Load() + q.onFlyRunningTasks.Load())
}
