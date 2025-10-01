package abstract

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maxbolgarin/lang"
)

// Result represents the outcome of a task execution.
type resultV2[T any] struct {
	Value T
	Err   error
}

// WorkerPool manages a pool of workers that process tasks concurrently.
type WorkerPoolV2[T any] struct {
	workers    int
	tasks      chan func() (T, error)
	results    chan resultV2[T]
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc

	started   atomic.Bool
	submitted atomic.Int64
	running   atomic.Int64
	finished  atomic.Int64
}

// NewWorkerPool creates a new worker pool with the specified number of workers and task queue capacity.
func NewWorkerPoolV2[T any](workers, queueCapacity int) *WorkerPoolV2[T] {
	if workers <= 0 {
		workers = 1
	}
	if queueCapacity <= 0 {
		queueCapacity = workers * 100
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPoolV2[T]{
		workers:    workers,
		tasks:      make(chan func() (T, error), queueCapacity),
		results:    make(chan resultV2[T], queueCapacity),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Start launches the worker goroutines.
func (p *WorkerPoolV2[T]) Start() {
	if p.started.Load() {
		return
	}

	p.wg.Add(p.workers)
	for range p.workers {
		lang.Go(nil, p.worker)
	}
	p.started.Store(true)
}

// Stop signals all workers to stop after completing their current tasks.
// It does not wait for them to complete.
func (p *WorkerPoolV2[T]) Stop() {
	if !p.started.Load() {
		return
	}
	p.cancelFunc()
	p.started.Store(false)
}

// worker is the goroutine that processes tasks.
func (p *WorkerPoolV2[T]) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			p.running.Add(1)
			value, err := task()
			select {
			case p.results <- resultV2[T]{Value: value, Err: err}:
				p.running.Add(-1)
				p.finished.Add(1)

			case <-p.ctx.Done():
				return
			}
		}
	}
}

// Submit adds a task to the pool and returns true if the task was accepted.
// Returns false if the pool is stopped or the task queue is full and the timeout is reached.
func (p *WorkerPoolV2[T]) Submit(task func() (T, error), timeoutRaw ...time.Duration) bool {
	if task == nil {
		return false
	}
	if p.IsStopped() {
		return false
	}

	if len(timeoutRaw) > 0 {
		timer := time.NewTimer(timeoutRaw[0])
		defer timer.Stop()

		select {
		case p.tasks <- task:
			p.submitted.Add(1)
			return true
		case <-timer.C:
			return false
		case <-p.ctx.Done():
			return false
		}
	}
	select {
	case p.tasks <- task:
		p.submitted.Add(1)
		return true
	case <-p.ctx.Done():
		return false
	}
}

// FetchResults fetches results from the pool.
// It returns when the number of results is equal to the number of submitted tasks AT THE TIME OF CALL!
// If the timeout is reached before the number of results is equal to the number of submitted tasks, it returns the results and errors.
// If some tasks are added after the call to FetchResults, they will not be fetched by this method (use FetchAllResults instead)
func (p *WorkerPoolV2[T]) FetchResults(timeoutRaw ...time.Duration) ([]T, []error) {
	var timeout time.Duration = time.Hour * 24 * 365
	if len(timeoutRaw) > 0 {
		timeout = timeoutRaw[0]
	}

	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()

	// Capture the count before the loop to avoid race condition
	expectedCount := int(p.submitted.Load())

	results := make([]T, 0, expectedCount)
	var errors []error

	for range expectedCount {
		select {
		case result := <-p.results:
			results = append(results, result.Value)
			errors = append(errors, result.Err)
			p.submitted.Add(-1)
			p.finished.Add(-1)
		case <-ctx.Done():
			return results, errors
		}
	}

	return results, errors
}

// FetchAllResults fetches all results from the pool.
// It returns when the number of results is equal to the number of submitted tasks!
// If the timeout is reached before the number of results is equal to the number of submitted tasks, it returns fetched results and errors.
// If some tasks are added after the call to FetchAllResults, they will be fetched by this method
func (p *WorkerPoolV2[T]) FetchAllResults(timeoutRaw ...time.Duration) ([]T, []error) {
	var timeout time.Duration = time.Hour * 24 * 365
	if len(timeoutRaw) > 0 {
		timeout = timeoutRaw[0]
	}

	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()

	results := make([]T, 0, p.submitted.Load())
	var errors []error

	for {
		expectedCount := int(p.submitted.Load())
		if expectedCount == 0 {
			return results, errors
		}

		select {
		case <-ctx.Done():
			return results, errors
		default:
		}

		resultsNow, errorsNow := p.FetchResults(timeoutRaw...)
		results = append(results, resultsNow...)
		errors = append(errors, errorsNow...)
	}
}

// Submitted returns the number of submitted tasks.
func (p *WorkerPoolV2[T]) Submitted() int {
	return int(p.submitted.Load())
}

// Running returns the number of running worker goroutines.
func (p *WorkerPoolV2[T]) Running() int {
	return int(p.running.Load())
}

// Finished returns the number of finished tasks.
func (p *WorkerPoolV2[T]) Finished() int {
	return int(p.finished.Load())
}

// IsStopped returns true if the worker pool has been stopped.
func (p *WorkerPoolV2[T]) IsStopped() bool {
	return !p.started.Load()
}
