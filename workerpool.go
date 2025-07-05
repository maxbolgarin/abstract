package abstract

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a function that can be executed by workers in the pool.
type Task func() (any, error)

// Result represents the outcome of a task execution.
type Result struct {
	Value any
	Err   error
}

// WorkerPool manages a pool of workers that process tasks concurrently.
type WorkerPool struct {
	workers    int
	tasks      chan Task
	results    chan Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	started    atomic.Bool
}

// NewWorkerPool creates a new worker pool with the specified number of workers and task queue capacity.
func NewWorkerPool(workers, queueCapacity int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	if queueCapacity <= 0 {
		queueCapacity = 100
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:    workers,
		tasks:      make(chan Task, queueCapacity),
		results:    make(chan Result, queueCapacity),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Start launches the worker goroutines.
func (p *WorkerPool) Start() {
	if p.started.Load() {
		return
	}

	p.wg.Add(p.workers)
	for range p.workers {
		go p.worker()
	}
	p.started.Store(true)
}

// worker is the goroutine that processes tasks.
func (p *WorkerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			value, err := task()
			select {
			case p.results <- Result{Value: value, Err: err}:
			case <-p.ctx.Done():
				return
			}
		}
	}
}

// Submit adds a task to the pool and returns true if the task was accepted.
// Returns false if the pool is stopped or the task queue is full and the timeout is reached.
func (p *WorkerPool) Submit(task Task, timeout time.Duration) bool {
	if task == nil {
		return false
	}
	if p.IsStopped() {
		return false
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case p.tasks <- task:
		return true
	case <-timer.C:
		return false
	case <-p.ctx.Done():
		return false
	}
}

// SubmitWait adds a task to the pool and waits for its completion, returning the result.
// If the timeout is reached before the task can be submitted or completed, it returns an error.
func (p *WorkerPool) SubmitWait(task Task, timeout time.Duration) (any, error) {
	if task == nil {
		return nil, errors.New("nil task submitted")
	}

	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()

	// Submit the task
	select {
	case p.tasks <- task:
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, errors.New("timeout submitting task")
		}
		return nil, ctx.Err()
	}

	// Wait for the result
	select {
	case result := <-p.results:
		return result.Value, result.Err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, errors.New("timeout waiting for result")
		}
		return nil, ctx.Err()
	}
}

// Results returns the channel that receives results from completed tasks.
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

// Stop signals all workers to stop after completing their current tasks.
// It does not wait for them to complete.
func (p *WorkerPool) Stop() {
	if !p.started.Load() {
		return
	}

	p.cancelFunc()
	p.started.Store(false)
}

// StopAndWait stops the worker pool and waits for all workers to complete.
// It returns true if workers completed within the timeout, false otherwise.
func (p *WorkerPool) StopAndWait(timeout time.Duration) bool {
	p.Stop()

	c := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Wait blocks until all workers have completed their tasks.
// This should only be called after Stop() or when all tasks have been submitted.
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// RunningWorkers returns the number of worker goroutines.
func (p *WorkerPool) RunningWorkers() int {
	return p.workers
}

// IsStopped returns true if the worker pool has been stopped.
func (p *WorkerPool) IsStopped() bool {
	return !p.started.Load()
}

// SafeWorkerPool is a thread-safe variant of WorkerPool.
type SafeWorkerPool struct {
	*WorkerPool
	mu sync.RWMutex
}

// NewSafeWorkerPool creates a new SafeWorkerPool.
func NewSafeWorkerPool(workers, queueCapacity int) *SafeWorkerPool {
	return &SafeWorkerPool{
		WorkerPool: NewWorkerPool(workers, queueCapacity),
	}
}

// Start launches the worker goroutines in a thread-safe manner.
func (p *SafeWorkerPool) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.WorkerPool.Start()
}

// Submit adds a task to the pool in a thread-safe manner.
func (p *SafeWorkerPool) Submit(task Task, timeout time.Duration) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.WorkerPool.Submit(task, timeout)
}

// SubmitWait adds a task to the pool and waits for its completion in a thread-safe manner.
func (p *SafeWorkerPool) SubmitWait(task Task, timeout time.Duration) (any, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.WorkerPool.SubmitWait(task, timeout)
}

// Stop signals all workers to stop in a thread-safe manner.
func (p *SafeWorkerPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.WorkerPool.Stop()
}

// StopAndWait stops the worker pool and waits for all workers to complete in a thread-safe manner.
func (p *SafeWorkerPool) StopAndWait(timeout time.Duration) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.WorkerPool.StopAndWait(timeout)
}

// IsStopped returns true if the worker pool has been stopped in a thread-safe manner.
func (p *SafeWorkerPool) IsStopped() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.WorkerPool.IsStopped()
}

// RunningWorkers returns the number of worker goroutines in a thread-safe manner.
func (p *SafeWorkerPool) RunningWorkers() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.WorkerPool.RunningWorkers()
}
