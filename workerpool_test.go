package abstract_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

func TestWorkerPoolBasicFunctionality(t *testing.T) {
	// Create a worker pool with 5 workers and a queue capacity of 10
	pool := abstract.NewWorkerPool(5, 10)
	pool.Start()
	defer pool.Stop()

	// Test submitting a simple task
	taskCompleted := false
	task := func() (interface{}, error) {
		taskCompleted = true
		return "success", nil
	}

	if !pool.Submit(task, time.Second) {
		t.Errorf("Failed to submit task to worker pool")
	}

	// Wait for task to complete
	result := <-pool.Results()
	if result.Err != nil {
		t.Errorf("Expected nil error, got %v", result.Err)
	}
	if result.Value != "success" {
		t.Errorf("Expected 'success', got %v", result.Value)
	}
	if !taskCompleted {
		t.Error("Task was not executed")
	}
}

func TestWorkerPoolSubmitWait(t *testing.T) {
	pool := abstract.NewWorkerPool(3, 5)
	pool.Start()
	defer pool.Stop()

	// Test SubmitWait with successful task
	value, err := pool.SubmitWait(func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		return 42, nil
	}, time.Second)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if value != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}

	// Test SubmitWait with error task
	expectedErr := errors.New("task error")
	value, err = pool.SubmitWait(func() (interface{}, error) {
		return nil, expectedErr
	}, time.Second)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if value != nil {
		t.Errorf("Expected nil value, got %v", value)
	}

	// Test SubmitWait with nil task
	_, err = pool.SubmitWait(nil, time.Second)
	if err == nil {
		t.Error("Expected error for nil task, got nil")
	}
}

func TestWorkerPoolConcurrency(t *testing.T) {
	// Number of workers and tasks
	workerCount := 10
	taskCount := 100

	pool := abstract.NewWorkerPool(workerCount, taskCount)
	pool.Start()
	defer pool.Stop()

	var counter int32
	var wg sync.WaitGroup
	wg.Add(taskCount)

	// Process results in a separate goroutine
	go func() {
		for i := 0; i < taskCount; i++ {
			<-pool.Results()
			wg.Done()
		}
	}()

	// Submit tasks
	for i := 0; i < taskCount; i++ {
		if !pool.Submit(func() (interface{}, error) {
			atomic.AddInt32(&counter, 1)
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		}, time.Second) {
			t.Errorf("Failed to submit task %d", i)
		}
	}

	// Wait for all tasks to complete
	wg.Wait()

	// Check if all tasks were executed
	if int(counter) != taskCount {
		t.Errorf("Expected %d tasks to be executed, got %d", taskCount, counter)
	}
}

func TestWorkerPoolTimeout(t *testing.T) {
	// Create a worker pool with a small capacity
	pool := abstract.NewWorkerPool(1, 1)
	pool.Start()

	// Fill the queue with slow tasks
	for i := 0; i < 2; i++ {
		pool.Submit(func() (any, error) {
			time.Sleep(500 * time.Millisecond)
			return nil, nil
		}, time.Second)
	}

	// This task should time out
	result := pool.Submit(func() (any, error) {
		return "should not execute", nil
	}, 50*time.Millisecond)

	if result {
		t.Error("Expected task submission to time out, but it succeeded")
	}

	// Cleanup
	pool.Stop()
}

func TestWorkerPoolShutdown(t *testing.T) {
	pool := abstract.NewWorkerPool(3, 10)
	pool.Start()

	// Submit some tasks
	for i := 0; i < 5; i++ {
		pool.Submit(func() (any, error) {
			time.Sleep(100 * time.Millisecond)
			return nil, nil
		}, time.Second)
	}

	// Stop and wait for completion
	completed := pool.StopAndWait(time.Second)
	if !completed {
		t.Error("Expected pool to stop within timeout")
	}

	if !pool.IsStopped() {
		t.Error("Pool should be marked as stopped")
	}

	// Submitting after stop should fail
	if pool.Submit(func() (any, error) {
		return nil, nil
	}, time.Second) {
		t.Error("Should not be able to submit tasks after stop")
	}
}

func TestSafeWorkerPool(t *testing.T) {
	pool := abstract.NewSafeWorkerPool(3, 10)
	pool.Start()
	defer pool.Stop()

	// Test concurrent operations on safe pool
	var wg sync.WaitGroup
	wg.Add(4)

	// Goroutine 1: Submit tasks
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			pool.Submit(func() (any, error) {
				time.Sleep(time.Millisecond)
				return i, nil
			}, time.Second)
		}
	}()

	// Goroutine 2: Submit and wait for tasks
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			_, err := pool.SubmitWait(func() (any, error) {
				return i, nil
			}, 2*time.Second)
			if err != nil {
				t.Errorf("SubmitWait failed: %v", err)
			}
		}
	}()

	// Goroutine 3: Check status
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			workers := pool.RunningWorkers()
			if workers != 3 {
				t.Errorf("Expected 3 workers, got %d", workers)
			}
			stopped := pool.IsStopped()
			if stopped {
				t.Error("Pool should not be stopped")
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Goroutine 4: Collect results
	go func() {
		defer wg.Done()
		count := 0
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) && count < 15 {
			select {
			case <-pool.Results():
				count++
			case <-time.After(50 * time.Millisecond):
				// Just wait a bit
			}
		}
	}()

	wg.Wait()
}

func TestWorkerPoolEdgeCases(t *testing.T) {
	// Test with zero or negative parameters
	pool := abstract.NewWorkerPool(0, -1)
	if pool.RunningWorkers() != 1 {
		t.Errorf("Expected 1 worker, got %d", pool.RunningWorkers())
	}

	// Test calling Start multiple times
	pool.Start()
	workers := pool.RunningWorkers()
	pool.Start() // Should be a no-op
	if pool.RunningWorkers() != workers {
		t.Error("Calling Start multiple times should not change worker count")
	}

	// Test nil task submission
	if pool.Submit(nil, time.Second) {
		t.Error("Submitting nil task should return false")
	}

	// Test stopping twice
	pool.Stop()
	wasStarted := !pool.IsStopped()
	pool.Stop() // Should be a no-op
	if pool.IsStopped() != !wasStarted {
		t.Error("Calling Stop multiple times should not change state")
	}

	// Test Wait after Stop
	pool = abstract.NewWorkerPool(2, 5)
	pool.Start()
	pool.Stop()

	// Wait should return immediately if pool is already stopped
	doneCh := make(chan struct{})
	go func() {
		pool.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// Expected behavior
	case <-time.After(500 * time.Millisecond):
		t.Error("Wait should return quickly after Stop")
	}
}
