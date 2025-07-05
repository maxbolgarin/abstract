package abstract

import (
	"context"
	"sync"
	"time"

	"github.com/maxbolgarin/lang"
)

// StartUpdater starts a panic-safe goroutine that executes a function periodically
// at the specified interval. The goroutine gracefully stops when the context is canceled.
//
// This function is useful for creating background tasks that need to run at regular
// intervals, such as health checks, data cleanup, or periodic synchronization.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - interval: Duration between function executions
//   - l: Logger for panic recovery and error handling
//   - f: Function to execute periodically
//
// Example usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	StartUpdater(ctx, 30*time.Second, logger, func() {
//		fmt.Println("Performing periodic health check...")
//		// Health check logic here
//	})
//
//	// Updater will stop when context is canceled
//	time.Sleep(5*time.Minute)
//	cancel() // Gracefully stops the updater
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

// StartUpdaterNow starts a panic-safe goroutine that executes a function immediately
// and then continues to execute it periodically at the specified interval.
//
// This is useful when you want to ensure the first execution happens right away
// instead of waiting for the first interval to pass.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - interval: Duration between subsequent function executions
//   - l: Logger for panic recovery and error handling
//   - f: Function to execute immediately and then periodically
//
// Example usage:
//
//	StartUpdaterNow(ctx, 1*time.Minute, logger, func() {
//		fmt.Println("Syncing data...") // Executes immediately
//		// Data synchronization logic
//	})
//	// Function runs immediately, then every minute thereafter
func StartUpdaterNow(ctx context.Context, interval time.Duration, l lang.Logger, f func()) {
	f()
	StartUpdater(ctx, interval, l, f)
}

// StartUpdaterWithShutdown starts a panic-safe goroutine that executes a function
// periodically and runs a shutdown function when the context is canceled.
//
// This pattern is useful when you need to perform cleanup operations when
// the periodic task is stopped.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - interval: Duration between function executions
//   - l: Logger for panic recovery and error handling
//   - f: Function to execute periodically
//   - shutdown: Function to execute when context is canceled
//
// Example usage:
//
//	StartUpdaterWithShutdown(ctx, 10*time.Second, logger,
//		func() {
//			// Periodic work
//			processQueue()
//		},
//		func() {
//			// Cleanup when stopping
//			fmt.Println("Shutting down queue processor...")
//			saveState()
//		},
//	)
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

// StartUpdaterWithShutdownChan starts a panic-safe goroutine that executes a function
// periodically and stops when either the context is canceled or the shutdown channel
// receives a signal.
//
// This provides an additional way to stop the updater without canceling the context,
// which can be useful in complex applications where you want fine-grained control.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - interval: Duration between function executions
//   - l: Logger for panic recovery and error handling
//   - c: Channel that signals shutdown when closed or receives a value
//   - f: Function to execute periodically
//
// Example usage:
//
//	shutdown := make(chan struct{})
//
//	StartUpdaterWithShutdownChan(ctx, 5*time.Second, logger, shutdown, func() {
//		fmt.Println("Periodic task running...")
//	})
//
//	// Later, signal shutdown
//	close(shutdown) // Stops the updater
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

// StartCycle starts a panic-safe goroutine that continuously executes a function
// in a tight loop until the context is canceled.
//
// Warning: This creates a high-CPU usage pattern. Use with caution and ensure
// the function either includes appropriate delays or processes work efficiently.
// Consider using StartUpdater with a small interval instead if periodic execution
// is acceptable.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - l: Logger for panic recovery and error handling
//   - f: Function to execute continuously
//
// Example usage:
//
//	StartCycle(ctx, logger, func() {
//		// Process available work or include a small delay
//		if work := getNextWork(); work != nil {
//			processWork(work)
//		} else {
//			time.Sleep(10 * time.Millisecond) // Prevent 100% CPU usage
//		}
//	})
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

// StartCycleWithShutdown starts a panic-safe goroutine that continuously executes
// a function until either the context is canceled or the shutdown channel is signaled.
//
// This provides additional control over stopping the cycle independently of context
// cancellation.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - l: Logger for panic recovery and error handling
//   - shutdown: Channel that signals shutdown when closed or receives a value
//   - f: Function to execute continuously
//
// Example usage:
//
//	shutdown := make(chan struct{})
//
//	StartCycleWithShutdown(ctx, logger, shutdown, func() {
//		// Continuous processing with controlled shutdown
//		processNextItem()
//	})
//
//	// Stop the cycle when needed
//	close(shutdown)
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

// StartCycleWithChan starts a panic-safe goroutine that processes values from a channel,
// executing the provided function for each received value until the context is canceled
// or the channel is closed.
//
// This is useful for implementing worker patterns where you need to process items
// from a queue or channel continuously.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - l: Logger for panic recovery and error handling
//   - c: Channel to receive values from
//   - f: Function to execute for each received value
//
// Example usage:
//
//	workChan := make(chan WorkItem, 100)
//
//	StartCycleWithChan(ctx, logger, workChan, func(item WorkItem) {
//		fmt.Printf("Processing work item: %v\n", item)
//		// Process the work item
//		item.Process()
//	})
//
//	// Send work to be processed
//	workChan <- WorkItem{ID: "task1", Data: "some data"}
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

// StartCycleWithChanAndShutdown starts a panic-safe goroutine that processes values
// from a channel until the context is canceled, the channel is closed, or a shutdown
// signal is received.
//
// This provides the most control over the lifecycle of a channel processing goroutine.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - l: Logger for panic recovery and error handling
//   - c: Channel to receive values from
//   - shutdown: Channel that signals shutdown when closed or receives a value
//   - f: Function to execute for each received value
//
// Example usage:
//
//	workChan := make(chan Task, 50)
//	shutdown := make(chan struct{})
//
//	StartCycleWithChanAndShutdown(ctx, logger, workChan, shutdown, func(task Task) {
//		fmt.Printf("Processing task: %v\n", task.ID)
//		task.Execute()
//	})
//
//	// Graceful shutdown
//	close(shutdown)
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

// RateProcessor manages a pool of workers to process tasks with rate limiting.
// It ensures that tasks are processed at a controlled rate, preventing system
// overload while maintaining high throughput.
//
// The rate processor is useful for scenarios where you need to:
//   - Limit API calls per second
//   - Control database query rates
//   - Manage external service interactions
//   - Prevent overwhelming downstream systems
//
// Example usage:
//
//	ctx := context.Background()
//	processor := NewRateProcessor(ctx, 10) // Max 10 tasks per second
//
//	// Add tasks
//	for i := 0; i < 100; i++ {
//		task := func(ctx context.Context) error {
//			// API call or other rate-limited operation
//			return makeAPICall(ctx, i)
//		}
//		processor.AddTask(task)
//	}
//
//	// Wait for completion and get any errors
//	errors := processor.Wait()
//	if len(errors) > 0 {
//		fmt.Printf("Encountered %d errors during processing\n", len(errors))
//	}
type RateProcessor struct {
	tasks   chan func(context.Context) error
	limiter <-chan time.Time
	wg      sync.WaitGroup
	errs    *SafeSlice[error]
}

// NewRateProcessor creates and starts a new RateProcessor with the specified
// maximum tasks per second rate limit.
//
// Parameters:
//   - ctx: Context for controlling the lifecycle of worker goroutines
//   - maxPerSecond: Maximum number of tasks to process per second
//
// Returns:
//   - A configured and started RateProcessor ready to accept tasks
//
// Example usage:
//
//	// Create a processor that handles max 5 tasks per second
//	processor := NewRateProcessor(ctx, 5)
//	defer processor.Wait() // Ensure cleanup
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
// The task will be executed by one of the workers when a rate limit slot
// becomes available.
//
// Note: This method will block if the task queue is full. Consider using
// a separate goroutine or adding timeout logic if non-blocking behavior
// is required.
//
// Parameters:
//   - task: Function to execute with rate limiting applied
//
// Example usage:
//
//	processor.AddTask(func(ctx context.Context) error {
//		response, err := http.Get("https://api.example.com/data")
//		if err != nil {
//			return fmt.Errorf("API call failed: %w", err)
//		}
//		defer response.Body.Close()
//
//		// Process response...
//		return nil
//	})
func (p *RateProcessor) AddTask(task func(context.Context) error) {
	p.tasks <- task
}

// Wait closes the task queue and waits for all workers to complete their
// current tasks. It returns all errors that occurred during task execution.
//
// This method should be called when no more tasks will be added to ensure
// proper cleanup and to retrieve any errors that occurred.
//
// Returns:
//   - A slice of all errors that occurred during task processing
//
// Example usage:
//
//	// Add all tasks...
//	for _, task := range tasks {
//		processor.AddTask(task)
//	}
//
//	// Wait for completion and handle errors
//	if errors := processor.Wait(); len(errors) > 0 {
//		for i, err := range errors {
//			log.Printf("Task error %d: %v", i+1, err)
//		}
//	}
func (p *RateProcessor) Wait() []error {
	close(p.tasks)
	p.wg.Wait()
	return p.errs.Copy()
}

// worker is the internal goroutine function that processes tasks with rate limiting.
// Each worker waits for the rate limiter before processing the next task.
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
