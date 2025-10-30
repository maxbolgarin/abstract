package abstract

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maxbolgarin/lang"
)

// result represents the outcome of a task execution.
type result[T any] struct {
	res T
	err error
}

// WorkerPool manages a pool of workers that process context-aware tasks concurrently.
// It provides advanced metrics and graceful shutdown capabilities.
type WorkerPool[T any] struct {
	workers  int
	tasks    chan func(ctx context.Context) (T, error)
	results  chan result[T]
	stopChan chan struct{}
	wg       sync.WaitGroup

	logger lang.Logger

	isPoolStarted     atomic.Bool
	onFlyRunningTasks atomic.Int64
	tasksInQueue      atomic.Int64
	finishedTasks     atomic.Int64
	totalTasks        atomic.Int64
}

// NewWorkerPool creates a new context-aware worker pool with the specified number of workers and task queue capacity.
func NewWorkerPool[T any](workers int, queueCapacity int, logger ...lang.Logger) *WorkerPool[T] {
	if workers <= 0 {
		workers = 1
	}
	if queueCapacity <= 0 {
		queueCapacity = workers * 100
	}

	return &WorkerPool[T]{
		workers:  workers,
		tasks:    make(chan func(ctx context.Context) (T, error), queueCapacity),
		results:  make(chan result[T], queueCapacity),
		stopChan: make(chan struct{}),
		logger:   lang.First(logger),
	}
}

// Start launches the worker goroutines.
func (p *WorkerPool[T]) Start(ctx context.Context) {
	if !p.isPoolStarted.CompareAndSwap(false, true) {
		return
	}
	p.wg.Add(p.workers)
	for range p.workers {
		lang.Go(p.logger, func() {
			p.worker(ctx)
		})
	}
}

// Shutdown signals all workers to stop after completing their current tasks.
// It waits for all in-flight and queued tasks to complete or until the context is done.
func (p *WorkerPool[T]) Shutdown(ctx context.Context) error {
	if !p.isPoolStarted.CompareAndSwap(true, false) {
		return nil
	}
	close(p.stopChan)
	close(p.tasks)

	// Wait for all workers to finish
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
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
func (p *WorkerPool[T]) StopNoWait() {
	if !p.isPoolStarted.CompareAndSwap(true, false) {
		return
	}
	close(p.stopChan)
	close(p.tasks)
}

// Submit adds a task to the pool and returns true if the task was accepted.
// Returns false if the pool is stopped or the context is done.
func (p *WorkerPool[T]) Submit(ctx context.Context, task func(ctx context.Context) (T, error)) bool {
	if task == nil {
		return false
	}
	if !p.isPoolStarted.Load() {
		return false
	}

	select {
	case p.tasks <- task:
		p.totalTasks.Add(1)
		p.tasksInQueue.Add(1)
		return true

	case <-p.stopChan:
		return false

	case <-ctx.Done():
		return false
	}
}

// worker is the goroutine that processes tasks.
func (p *WorkerPool[T]) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				// Channel closed, drain remaining tasks
				return
			}
			p.tasksInQueue.Add(-1)

			p.onFlyRunningTasks.Add(1)
			value, err := task(ctx)
			p.onFlyRunningTasks.Add(-1)

			p.results <- result[T]{res: value, err: err}
			p.finishedTasks.Add(1)

		case <-ctx.Done():
			return
		case <-p.stopChan:
			// Stop signal received, but continue processing pending tasks
			// The tasks channel will be closed, causing the worker to exit after draining
		}
	}
}

// FetchResults fetches results from the pool.
// It returns when the number of results is equal to the number of finished tasks AT THE TIME OF CALL!
// If the context is done before all results are fetched, it returns the results and errors collected so far.
// If some tasks are added after the call to FetchResults, they will not be fetched by this method (use FetchAllResults instead).
func (p *WorkerPool[T]) FetchResults(ctx context.Context) ([]T, []error) {
	// Capture the count before the loop to avoid race condition
	expectedCount := int(p.finishedTasks.Load())

	results := make([]T, 0, expectedCount)
	errors := make([]error, 0, expectedCount)

	for range expectedCount {
		select {
		case result := <-p.results:
			results = append(results, result.res)
			errors = append(errors, result.err)
			p.finishedTasks.Add(-1)

		case <-ctx.Done():
			return results, errors
		}
	}

	return results, errors
}

// FetchAllResults fetches all results from the pool.
// It waits until all submitted tasks have finished and returns their results.
// If the context is done before all results are fetched, it returns fetched results and errors.
func (p *WorkerPool[T]) FetchAllResults(ctx context.Context) ([]T, []error) {
	results := make([]T, 0)
	errors := make([]error, 0)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		// Check if all tasks are done
		finished := int(p.finishedTasks.Load())
		if finished == 0 && p.tasksInQueue.Load() == 0 && p.onFlyRunningTasks.Load() == 0 {
			return results, errors
		}

		if finished > 0 {
			// Fetch available results
			resultsNow, errorsNow := p.FetchResults(ctx)
			results = append(results, resultsNow...)
			errors = append(errors, errorsNow...)
		}

		select {
		case <-ctx.Done():
			return results, errors
		case <-ticker.C:
			// Continue checking
		}
	}
}

// TasksInQueue returns the number of tasks in the queue.
func (p *WorkerPool[T]) TasksInQueue() int {
	return int(p.tasksInQueue.Load())
}

// OnFlyRunningTasks returns the number of currently executing tasks.
func (p *WorkerPool[T]) OnFlyRunningTasks() int {
	return int(p.onFlyRunningTasks.Load())
}

// FinishedTasks returns the number of finished tasks waiting to be fetched.
func (p *WorkerPool[T]) FinishedTasks() int {
	return int(p.finishedTasks.Load())
}

// TotalTasks returns the total number of tasks submitted to the pool.
func (p *WorkerPool[T]) TotalTasks() int {
	return int(p.totalTasks.Load())
}

// IsPoolStarted returns true if the worker pool has been started.
func (p *WorkerPool[T]) IsPoolStarted() bool {
	return p.isPoolStarted.Load()
}
