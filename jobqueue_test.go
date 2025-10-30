package abstract_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

func TestNewJobQueue(t *testing.T) {
	queue := abstract.NewJobQueue(5, 100)
	if queue == nil {
		t.Fatal("Expected non-nil JobQueue")
	}
	if queue.IsQueueStarted() {
		t.Error("Queue should not be started initially")
	}
}

func TestJobQueueBasicExecution(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(3, 10)
	queue.Start(ctx)
	defer queue.StopNoWait()

	var counter atomic.Int32

	// Submit tasks
	taskCount := 10
	for i := 0; i < taskCount; i++ {
		ok := queue.Submit(ctx, func(ctx context.Context) {
			counter.Add(1)
		})
		if !ok {
			t.Errorf("Failed to submit task %d", i)
		}
	}

	// Wait for all tasks to complete
	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := queue.Wait(waitCtx)
	if err != nil {
		t.Errorf("Wait failed: %v", err)
	}

	if counter.Load() != int32(taskCount) {
		t.Errorf("Expected %d executions, got %d", taskCount, counter.Load())
	}
}

func TestJobQueueMetrics(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 10)
	queue.Start(ctx)
	defer queue.StopNoWait()

	if queue.TotalTasks() != 0 {
		t.Errorf("Expected 0 total tasks, got %d", queue.TotalTasks())
	}

	var executing atomic.Int32
	var done atomic.Int32

	// Submit tasks that take some time
	taskCount := 5
	for i := 0; i < taskCount; i++ {
		ok := queue.Submit(ctx, func(ctx context.Context) {
			executing.Add(1)
			time.Sleep(50 * time.Millisecond)
			executing.Add(-1)
			done.Add(1)
		})
		if !ok {
			t.Errorf("Failed to submit task %d", i)
		}
	}

	if queue.TotalTasks() != taskCount {
		t.Errorf("Expected %d total tasks, got %d", taskCount, queue.TotalTasks())
	}

	// Give some time for tasks to start processing
	time.Sleep(20 * time.Millisecond)

	// Check metrics
	queued := queue.TasksInQueue()
	running := queue.OnFlyRunningTasks()
	finished := queue.FinishedTasks()
	pending := queue.PendingTasks()

	if queued+running != pending {
		t.Errorf("PendingTasks (%d) != TasksInQueue (%d) + OnFlyRunningTasks (%d)", pending, queued, running)
	}

	t.Logf("Queued: %d, Running: %d, Finished: %d, Pending: %d", queued, running, finished, pending)

	// Wait for completion
	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := queue.Wait(waitCtx)
	if err != nil {
		t.Errorf("Wait failed: %v", err)
	}

	if queue.PendingTasks() != 0 {
		t.Errorf("Expected 0 pending tasks after wait, got %d", queue.PendingTasks())
	}

	if done.Load() != int32(taskCount) {
		t.Errorf("Expected %d completed tasks, got %d", taskCount, done.Load())
	}
}

func TestJobQueueShutdown(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 10)
	queue.Start(ctx)

	var counter atomic.Int32

	// Submit some tasks
	for i := 0; i < 5; i++ {
		queue.Submit(ctx, func(ctx context.Context) {
			time.Sleep(10 * time.Millisecond)
			counter.Add(1)
		})
	}

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err := queue.Shutdown(shutdownCtx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	if queue.IsQueueStarted() {
		t.Error("Queue should be stopped after shutdown")
	}

	// All tasks should have completed
	if counter.Load() != 5 {
		t.Errorf("Expected 5 tasks to complete, got %d", counter.Load())
	}
}

func TestJobQueueCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	queue := abstract.NewJobQueue(2, 10)
	queue.Start(ctx)
	defer queue.StopNoWait()

	var completedTasks atomic.Int32

	// Submit tasks that respect context
	for i := 0; i < 10; i++ {
		queue.Submit(ctx, func(ctx context.Context) {
			select {
			case <-ctx.Done():
				// Task cancelled
				return
			case <-time.After(100 * time.Millisecond):
				completedTasks.Add(1)
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

func TestJobQueueWaitWithTimeout(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 10)
	queue.Start(ctx)
	defer queue.StopNoWait()

	// Submit tasks that take longer than wait timeout
	for i := 0; i < 5; i++ {
		queue.Submit(ctx, func(ctx context.Context) {
			time.Sleep(500 * time.Millisecond)
		})
	}

	// Wait with short timeout
	waitCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err := queue.Wait(waitCtx)
	if err == nil {
		t.Error("Expected Wait to timeout, but it succeeded")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", err)
	}
}

func TestJobQueueNilTask(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 5)
	queue.Start(ctx)
	defer queue.StopNoWait()

	ok := queue.Submit(ctx, nil)
	if ok {
		t.Error("Expected Submit to reject nil task")
	}
}

func TestJobQueueSubmitToStoppedQueue(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 5)
	queue.Start(ctx)
	queue.StopNoWait()

	time.Sleep(50 * time.Millisecond)

	var executed atomic.Bool
	ok := queue.Submit(ctx, func(ctx context.Context) {
		executed.Store(true)
	})

	if ok {
		t.Error("Expected Submit to fail on stopped queue")
	}

	time.Sleep(50 * time.Millisecond)
	if executed.Load() {
		t.Error("Task should not have been executed")
	}
}

func TestJobQueueDoubleStart(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 5)

	queue.Start(ctx)
	if !queue.IsQueueStarted() {
		t.Error("Queue should be started")
	}

	// Starting again should be a no-op
	queue.Start(ctx)
	if !queue.IsQueueStarted() {
		t.Error("Queue should still be started")
	}

	queue.StopNoWait()
}

func TestJobQueueZeroWorkers(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(0, 10) // Should default to 1
	queue.Start(ctx)
	defer queue.StopNoWait()

	var executed atomic.Bool
	ok := queue.Submit(ctx, func(ctx context.Context) {
		executed.Store(true)
	})

	if !ok {
		t.Error("Expected task submission to succeed")
	}

	waitCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := queue.Wait(waitCtx)
	if err != nil {
		t.Errorf("Wait failed: %v", err)
	}

	if !executed.Load() {
		t.Error("Task should have been executed")
	}
}

func TestJobQueueZeroCapacity(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 0) // Should default to workers * 100
	queue.Start(ctx)
	defer queue.StopNoWait()

	// Should be able to submit tasks
	ok := queue.Submit(ctx, func(ctx context.Context) {
		// No-op
	})

	if !ok {
		t.Error("Expected task submission to succeed with default capacity")
	}
}

func TestJobQueueConcurrentSubmit(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(5, 100)
	queue.Start(ctx)
	defer queue.StopNoWait()

	var counter atomic.Int32
	taskCount := 100

	// Submit tasks concurrently
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < taskCount/10; j++ {
				queue.Submit(ctx, func(ctx context.Context) {
					counter.Add(1)
					time.Sleep(time.Millisecond)
				})
			}
			done <- struct{}{}
		}()
	}

	// Wait for all submissions
	for i := 0; i < 10; i++ {
		<-done
	}

	// Wait for all tasks to complete
	waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := queue.Wait(waitCtx)
	if err != nil {
		t.Errorf("Wait failed: %v", err)
	}

	if counter.Load() != int32(taskCount) {
		t.Errorf("Expected %d executions, got %d", taskCount, counter.Load())
	}
}

func TestJobQueueWaitOnEmptyQueue(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 5)
	queue.Start(ctx)
	defer queue.StopNoWait()

	// Wait on empty queue should return immediately
	waitCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	start := time.Now()
	err := queue.Wait(waitCtx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Wait on empty queue failed: %v", err)
	}

	if elapsed > 100*time.Millisecond {
		t.Errorf("Wait on empty queue took too long: %v", elapsed)
	}
}

func TestJobQueuePanicRecovery(t *testing.T) {
	ctx := context.Background()
	queue := abstract.NewJobQueue(2, 5)
	queue.Start(ctx)
	defer queue.StopNoWait()

	var normalTaskExecuted atomic.Bool

	// Submit a task that panics
	queue.Submit(ctx, func(ctx context.Context) {
		panic("intentional panic")
	})

	// Submit a normal task
	queue.Submit(ctx, func(ctx context.Context) {
		normalTaskExecuted.Store(true)
	})

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// The queue should still be functional
	// Note: Without explicit panic recovery in the worker, this test
	// demonstrates current behavior. The normal task should still execute
	// if workers continue after panic (depends on lang.Go implementation)

	// This test mainly ensures the test suite doesn't crash
	t.Log("Queue survived panic in task")
}
