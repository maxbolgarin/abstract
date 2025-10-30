package abstract_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

func TestNewWorkerPool(t *testing.T) {
	pool := abstract.NewWorkerPool[int](5, 100)
	if pool == nil {
		t.Fatal("Expected non-nil WorkerPool")
	}
	if pool.IsPoolStarted() {
		t.Error("Pool should not be started initially")
	}
}

func TestWorkerPoolBasicExecution(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](3, 10)
	pool.Start(ctx)
	defer pool.StopNoWait()

	// Submit tasks
	for i := 0; i < 5; i++ {
		val := i
		ok := pool.Submit(ctx, func(ctx context.Context) (int, error) {
			return val * 2, nil
		})
		if !ok {
			t.Errorf("Failed to submit task %d", i)
		}
	}

	// Wait a bit for tasks to complete
	time.Sleep(100 * time.Millisecond)

	// Fetch results
	results, errs := pool.FetchResults(ctx)
	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	// Check for errors
	for i, err := range errs {
		if err != nil {
			t.Errorf("Task %d returned error: %v", i, err)
		}
	}

	// Verify results contain expected values
	resultMap := make(map[int]bool)
	for _, r := range results {
		resultMap[r] = true
	}
	for i := 0; i < 5; i++ {
		expected := i * 2
		if !resultMap[expected] {
			t.Errorf("Expected result %d not found", expected)
		}
	}
}

func TestWorkerPoolWithErrors(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[string](2, 5)
	pool.Start(ctx)
	defer pool.StopNoWait()

	expectedErr := errors.New("task error")

	// Submit tasks with errors
	for i := 0; i < 3; i++ {
		val := i
		ok := pool.Submit(ctx, func(ctx context.Context) (string, error) {
			if val%2 == 0 {
				return "", expectedErr
			}
			return "success", nil
		})
		if !ok {
			t.Errorf("Failed to submit task %d", i)
		}
	}

	time.Sleep(100 * time.Millisecond)

	results, errs := pool.FetchResults(ctx)
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
	if len(errs) != 3 {
		t.Errorf("Expected 3 error slots, got %d", len(errs))
	}

	// Count errors
	errorCount := 0
	for _, err := range errs {
		if err != nil {
			errorCount++
		}
	}
	if errorCount != 2 {
		t.Errorf("Expected 2 errors, got %d", errorCount)
	}
}

func TestWorkerPoolMetrics(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](2, 10)
	pool.Start(ctx)
	defer pool.StopNoWait()

	if pool.TotalTasks() != 0 {
		t.Errorf("Expected 0 total tasks, got %d", pool.TotalTasks())
	}

	// Submit tasks
	taskCount := 5
	for i := 0; i < taskCount; i++ {
		pool.Submit(ctx, func(ctx context.Context) (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 1, nil
		})
	}

	if pool.TotalTasks() != taskCount {
		t.Errorf("Expected %d total tasks, got %d", taskCount, pool.TotalTasks())
	}

	// Give some time for tasks to start processing
	time.Sleep(20 * time.Millisecond)

	// Check metrics
	queued := pool.TasksInQueue()
	running := pool.OnFlyRunningTasks()
	finished := pool.FinishedTasks()

	if queued+running+finished != taskCount {
		t.Logf("Queued: %d, Running: %d, Finished: %d", queued, running, finished)
		// Note: Due to timing, this might not always be exact, so we just log
	}

	// Wait for completion
	time.Sleep(300 * time.Millisecond)

	if pool.FinishedTasks() != taskCount {
		t.Errorf("Expected %d finished tasks, got %d", taskCount, pool.FinishedTasks())
	}
}

func TestWorkerPoolShutdown(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](2, 10)
	pool.Start(ctx)

	// Submit some tasks
	for i := 0; i < 5; i++ {
		pool.Submit(ctx, func(ctx context.Context) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 1, nil
		})
	}

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err := pool.Shutdown(shutdownCtx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	if pool.IsPoolStarted() {
		t.Error("Pool should be stopped after shutdown")
	}
}

func TestWorkerPoolCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := abstract.NewWorkerPool[int](2, 10)
	pool.Start(ctx)
	defer pool.StopNoWait()

	var completedTasks atomic.Int32

	// Submit tasks that respect context
	for i := 0; i < 10; i++ {
		pool.Submit(ctx, func(ctx context.Context) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				completedTasks.Add(1)
				return 1, nil
			}
		})
	}

	// Cancel context after a short time
	time.Sleep(50 * time.Millisecond)
	cancel()

	time.Sleep(200 * time.Millisecond)

	completed := completedTasks.Load()
	if completed >= 10 {
		t.Error("Expected some tasks to be cancelled, but all completed")
	}
	t.Logf("Completed tasks: %d out of 10", completed)
}

func TestWorkerPoolFetchAllResults(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](3, 10)
	pool.Start(ctx)
	defer pool.StopNoWait()

	// Submit all tasks with slight delay so they don't complete instantly
	for i := 0; i < 8; i++ {
		val := i
		pool.Submit(ctx, func(ctx context.Context) (int, error) {
			time.Sleep(20 * time.Millisecond)
			return val, nil
		})
	}

	// Fetch all results
	fetchCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	allResults, allErrors := pool.FetchAllResults(fetchCtx)

	if len(allResults) != 8 {
		t.Errorf("Expected 8 results, got %d", len(allResults))
	}

	for _, err := range allErrors {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

func TestWorkerPoolNilTask(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](2, 5)
	pool.Start(ctx)
	defer pool.StopNoWait()

	ok := pool.Submit(ctx, nil)
	if ok {
		t.Error("Expected Submit to reject nil task")
	}
}

func TestWorkerPoolSubmitToStoppedPool(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](2, 5)
	pool.Start(ctx)
	pool.StopNoWait()

	time.Sleep(50 * time.Millisecond)

	ok := pool.Submit(ctx, func(ctx context.Context) (int, error) {
		return 1, nil
	})

	if ok {
		t.Error("Expected Submit to fail on stopped pool")
	}
}

func TestWorkerPoolDoubleStart(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](2, 5)

	pool.Start(ctx)
	if !pool.IsPoolStarted() {
		t.Error("Pool should be started")
	}

	// Starting again should be a no-op
	pool.Start(ctx)
	if !pool.IsPoolStarted() {
		t.Error("Pool should still be started")
	}

	pool.StopNoWait()
}

func TestWorkerPoolZeroWorkers(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](0, 10) // Should default to 1
	pool.Start(ctx)
	defer pool.StopNoWait()

	ok := pool.Submit(ctx, func(ctx context.Context) (int, error) {
		return 42, nil
	})

	if !ok {
		t.Error("Expected task submission to succeed")
	}

	time.Sleep(100 * time.Millisecond)
	results, _ := pool.FetchResults(ctx)

	if len(results) != 1 || results[0] != 42 {
		t.Errorf("Expected result [42], got %v", results)
	}
}

func TestWorkerPoolZeroCapacity(t *testing.T) {
	ctx := context.Background()
	pool := abstract.NewWorkerPool[int](2, 0) // Should default to workers * 100
	pool.Start(ctx)
	defer pool.StopNoWait()

	// Should be able to submit tasks
	ok := pool.Submit(ctx, func(ctx context.Context) (int, error) {
		return 1, nil
	})

	if !ok {
		t.Error("Expected task submission to succeed with default capacity")
	}
}
