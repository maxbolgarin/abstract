package abstract

import (
	"context"
	"sync"
	"time"

	"github.com/maxbolgarin/lang"
)

// StartUpdater starts a new panicsafe goroutine.
// It runs the provided function one time per the interval.
// It stops the goroutine when the context is canceled.
func StartUpdater(ctx context.Context, interval time.Duration, l lang.Logger, f func()) {
	lang.Go(l, func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f()
			case <-ctx.Done():
				return
			}
		}
	})
}

// StartUpdaterNow starts a new panicsafe goroutine.
// It runs the provided function right now and then one time per the interval.
// It stops the goroutine when the context is canceled.
func StartUpdaterNow(ctx context.Context, interval time.Duration, l lang.Logger, f func()) {
	f()
	StartUpdater(ctx, interval, l, f)
}

// StartUpdaterWithShutdown starts a new panicsafe goroutine.
// It runs the provided function one time per the interval.
// It runs the shutdown function and stops the goroutine when the context is canceled.
func StartUpdaterWithShutdown(ctx context.Context, interval time.Duration, l lang.Logger, f func(), shutdown func()) {
	lang.Go(l, func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f()
			case <-ctx.Done():
				shutdown()
				return
			}
		}
	})
}

// StartUpdaterWithShutdownChan starts a new panicsafe goroutine.
// It runs the provided function one time per the interval.
// It stops the goroutine when the context is canceled or the provided channel is closed.
func StartUpdaterWithShutdownChan(ctx context.Context, interval time.Duration, l lang.Logger, c chan struct{}, f func()) {
	lang.Go(l, func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f()
			case <-c:
				return
			case <-ctx.Done():
				return
			}
		}
	})
}

// StartCycle starts a new panicsafe goroutine.
// It runs the provided function in a loop until the context is canceled.
func StartCycle(ctx context.Context, l lang.Logger, f func()) {
	lang.Go(l, func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				f()
			}
		}
	})
}

// StartCycleWithShutdown starts a new panicsafe goroutine.
// It runs the provided function in a loop until the context is canceled or the provided shutdown channel is closed.
func StartCycleWithShutdown(ctx context.Context, l lang.Logger, shutdown <-chan struct{}, f func()) {
	lang.Go(l, func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-shutdown:
				return
			default:
				f()
			}
		}
	})
}

// StartCycleWithChan starts a new panicsafe goroutine.
// It runs the provided function for each value from the provided channel.
// It stops the goroutine when the context is canceled or the provided channel is closed.
func StartCycleWithChan[T any](ctx context.Context, l lang.Logger, c <-chan T, f func(T)) {
	lang.Go(l, func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-c:
				f(val)
			}
		}
	})
}

// StartCycleWithChanAndShutdown starts a new panicsafe goroutine.
// It runs the provided function for each value from the provided channel.
// It stops the goroutine when the context is canceled or the provided shutdown channel is closed.
func StartCycleWithChanAndShutdown[T any](ctx context.Context, l lang.Logger, c <-chan T, shutdown <-chan struct{}, f func(T)) {
	lang.Go(l, func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-c:
				f(val)
			case <-shutdown:
				return
			}
		}
	})
}

// RateProcessor manages a pool of workers to process tasks with a rate limit.
type RateProcessor struct {
	tasks   chan func(context.Context) error
	limiter <-chan time.Time
	wg      sync.WaitGroup
	errs    *SafeSlice[error]
}

// NewRateProcessor initializes a new RateProcessor, that manages a pool of workers to process tasks with a rate limit.
func NewRateProcessor(ctx context.Context, maxPerSecond int) *RateProcessor {
	p := &RateProcessor{
		tasks:   make(chan func(context.Context) error, maxPerSecond),
		limiter: time.Tick(time.Second / time.Duration(maxPerSecond)),
		errs:    NewSafeSlice[error](),
	}

	for i := 0; i < maxPerSecond; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}

	return p
}

// AddTask adds a task to the worker pool's task queue.
func (p *RateProcessor) AddTask(task func(context.Context) error) {
	p.tasks <- task
}

// Wait closes down the worker pool and waits for all workers to complete.
// It returns a slice of errors that occurred during task execution.
func (p *RateProcessor) Wait() []error {
	close(p.tasks)
	p.wg.Wait()
	return p.errs.Copy()
}

func (p *RateProcessor) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			select {
			case <-p.limiter:
				if err := task(ctx); err != nil {
					p.errs.Append(err)
				}
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
