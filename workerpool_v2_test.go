package abstract_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

func TestWorkerPoolV2BasicFunctionality(t *testing.T) {
	// Create a worker pool with 5 workers and a queue capacity of 10
	pool := abstract.NewWorkerPoolV2[string](5, 10)
	pool.Start()
	defer pool.Stop()

	// Test submitting a simple task
	taskCompleted := false
	task := func() (string, error) {
		taskCompleted = true
		return "success", nil
	}

	if !pool.Submit(task) {
		t.Error("Failed to submit task to worker pool")
	}

	// Wait for task to complete
	results, errors := pool.FetchResults(time.Second)
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if len(errors) != 1 {
		t.Errorf("Expected 1 error slot, got %d", len(errors))
	}
	if errors[0] != nil {
		t.Errorf("Expected nil error, got %v", errors[0])
	}
	if results[0] != "success" {
		t.Errorf("Expected 'success', got %v", results[0])
	}
	if !taskCompleted {
		t.Error("Task was not executed")
	}
}

func TestWorkerPoolV2EdgeCases(t *testing.T) {
	// Test with zero or negative parameters
	pool := abstract.NewWorkerPoolV2[int](0, -1)

	// Test with negative workers
	pool = abstract.NewWorkerPoolV2[int](-5, 50)

	// Test calling Start multiple times
	pool.Start()
	pool.Start() // Should be a no-op

	// Test nil task submission
	if pool.Submit(nil) {
		t.Error("Submitting nil task should return false")
	}

	// Test IsStopped before and after start
	pool2 := abstract.NewWorkerPoolV2[int](2, 5)
	if !pool2.IsStopped() {
		t.Error("Pool should be stopped initially")
	}
	pool2.Start()
	if pool2.IsStopped() {
		t.Error("Pool should not be stopped after Start")
	}

	// Test stopping twice
	pool.Stop()
	wasStarted := !pool.IsStopped()
	pool.Stop() // Should be a no-op
	if pool.IsStopped() == wasStarted {
		t.Error("Pool should remain stopped after calling Stop multiple times")
	}

	// Test Stop when not started
	pool3 := abstract.NewWorkerPoolV2[int](2, 5)
	pool3.Stop() // Should not panic
	if !pool3.IsStopped() {
		t.Error("Pool should still be stopped after calling Stop without Start")
	}
}

func TestWorkerPoolV2SubmitWithTimeout(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](1, 1)
	pool.Start()
	defer pool.Stop()

	// Fill the queue with slow tasks
	submitted := pool.Submit(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	})
	if !submitted {
		t.Error("First task should be submitted successfully")
	}

	// Fill the buffer
	submitted = pool.Submit(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 2, nil
	})
	if !submitted {
		t.Error("Second task should be submitted successfully")
	}

	// This task should time out
	submitted = pool.Submit(func() (int, error) {
		return 3, nil
	}, 50*time.Millisecond)
	if submitted {
		t.Error("Expected task submission to time out, but it succeeded")
	}

	// Wait for tasks to complete
	time.Sleep(500 * time.Millisecond)
}

func TestWorkerPoolV2SubmitAfterStop(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 10)
	pool.Start()
	pool.Stop()

	// Submitting after stop should fail
	if pool.Submit(func() (int, error) {
		return 42, nil
	}) {
		t.Error("Should not be able to submit tasks after stop")
	}

	// Submitting with timeout after stop should also fail
	if pool.Submit(func() (int, error) {
		return 42, nil
	}, time.Second) {
		t.Error("Should not be able to submit tasks with timeout after stop")
	}
}

func TestWorkerPoolV2WithErrors(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 10)
	pool.Start()
	defer pool.Stop()

	expectedErr := errors.New("task error")

	// Submit task that returns error
	if !pool.Submit(func() (int, error) {
		return 0, expectedErr
	}) {
		t.Error("Failed to submit task")
	}

	// Submit task that succeeds
	if !pool.Submit(func() (int, error) {
		return 42, nil
	}) {
		t.Error("Failed to submit task")
	}

	results, errs := pool.FetchResults(time.Second)
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
	if len(errs) != 2 {
		t.Errorf("Expected 2 error slots, got %d", len(errs))
	}

	// Check that we got one error and one success
	foundError := false
	foundSuccess := false
	for i := range results {
		if errs[i] == expectedErr {
			foundError = true
			if results[i] != 0 {
				t.Errorf("Expected zero value for error case, got %v", results[i])
			}
		}
		if errs[i] == nil && results[i] == 42 {
			foundSuccess = true
		}
	}

	if !foundError {
		t.Error("Did not find expected error in results")
	}
	if !foundSuccess {
		t.Error("Did not find expected success in results")
	}
}

func TestWorkerPoolV2Concurrency(t *testing.T) {
	// Number of workers and tasks
	workerCount := 10
	taskCount := 100

	pool := abstract.NewWorkerPoolV2[int](workerCount, taskCount)
	pool.Start()
	defer pool.Stop()

	var counter int32

	// Submit tasks
	for i := 0; i < taskCount; i++ {
		if !pool.Submit(func() (int, error) {
			atomic.AddInt32(&counter, 1)
			time.Sleep(5 * time.Millisecond)
			return int(atomic.LoadInt32(&counter)), nil
		}) {
			t.Errorf("Failed to submit task %d", i)
		}
	}

	// Fetch all results
	results, errs := pool.FetchResults(5 * time.Second)

	if len(results) != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, len(results))
	}
	if len(errs) != taskCount {
		t.Errorf("Expected %d error slots, got %d", taskCount, len(errs))
	}

	// Check if all tasks were executed
	if int(counter) != taskCount {
		t.Errorf("Expected %d tasks to be executed, got %d", taskCount, counter)
	}

	// Check that all errors are nil
	for i, err := range errs {
		if err != nil {
			t.Errorf("Expected nil error at index %d, got %v", i, err)
		}
	}
}

func TestWorkerPoolV2FetchResultsTimeout(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()
	defer pool.Stop()

	// Submit some fast tasks
	for i := 0; i < 3; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 1, nil
		})
	}

	// Submit some slow tasks that won't complete in time
	for i := 0; i < 2; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(500 * time.Millisecond)
			return 2, nil
		})
	}

	// Fetch with short timeout - should get some but not all results
	results, errs := pool.FetchResults(200 * time.Millisecond)

	if len(results) == 0 {
		t.Error("Expected at least some results")
	}
	if len(results) != len(errs) {
		t.Errorf("Results and errors length mismatch: %d vs %d", len(results), len(errs))
	}
	if len(results) >= 5 {
		t.Error("Expected timeout to prevent all results from being fetched")
	}
}

func TestWorkerPoolV2FetchResultsNoTimeout(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[string](3, 10)
	pool.Start()
	defer pool.Stop()

	// Submit tasks
	taskCount := 5
	for i := 0; i < taskCount; i++ {
		index := i
		pool.Submit(func() (string, error) {
			time.Sleep(10 * time.Millisecond)
			return "task" + string(rune('0'+index)), nil
		})
	}

	// Fetch results without explicit timeout (uses default large timeout)
	results, errs := pool.FetchResults()

	if len(results) != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, len(results))
	}
	if len(errs) != taskCount {
		t.Errorf("Expected %d error slots, got %d", taskCount, len(errs))
	}

	// All errors should be nil
	for i, err := range errs {
		if err != nil {
			t.Errorf("Expected nil error at index %d, got %v", i, err)
		}
	}
}

func TestWorkerPoolV2StopDuringFetch(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()

	// Submit tasks
	for i := 0; i < 5; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		})
	}

	// Start fetching results in a goroutine
	doneCh := make(chan struct{})
	var results []int
	go func() {
		results, _ = pool.FetchResults(5 * time.Second)
		close(doneCh)
	}()

	// Stop the pool while fetching
	time.Sleep(150 * time.Millisecond)
	pool.Stop()

	// Wait for fetch to complete
	select {
	case <-doneCh:
		// Expected behavior - fetch should complete
		if len(results) == 0 {
			t.Error("Expected at least some results before stop")
		}
	case <-time.After(6 * time.Second):
		t.Error("FetchResults did not complete after stop")
	}
}

func TestWorkerPoolV2ConcurrentOperations(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](5, 50)
	pool.Start()
	defer pool.Stop()

	var wg sync.WaitGroup
	wg.Add(3)

	// Goroutine 1: Submit tasks
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			pool.Submit(func() (int, error) {
				time.Sleep(10 * time.Millisecond)
				return i, nil
			})
			time.Sleep(5 * time.Millisecond)
		}
	}()

	// Goroutine 2: Check status
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			stopped := pool.IsStopped()
			if stopped {
				t.Error("Pool should not be stopped")
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	// Goroutine 3: Fetch some results
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		results, _ := pool.FetchResults(500 * time.Millisecond)
		if len(results) == 0 {
			t.Error("Expected some results")
		}
	}()

	wg.Wait()
}

func TestWorkerPoolV2SubmitWithTimeoutSuccess(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](5, 10)
	pool.Start()
	defer pool.Stop()

	// Submit with timeout should succeed when queue has space
	submitted := pool.Submit(func() (int, error) {
		return 42, nil
	}, time.Second)

	if !submitted {
		t.Error("Task with timeout should be submitted successfully")
	}

	results, errs := pool.FetchResults(time.Second)
	if len(results) != 1 || results[0] != 42 {
		t.Error("Expected result 42")
	}
	if len(errs) != 1 || errs[0] != nil {
		t.Error("Expected no error")
	}
}

func TestWorkerPoolV2SubmitWhileStopping(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 5)
	pool.Start()

	// Submit some tasks
	for i := 0; i < 3; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		})
	}

	// Stop the pool
	pool.Stop()

	// Try to submit after stopping - should fail immediately
	submitted := pool.Submit(func() (int, error) {
		return 99, nil
	}, 100*time.Millisecond)

	if submitted {
		t.Error("Should not be able to submit after pool is stopped")
	}
}

func TestWorkerPoolV2ZeroResults(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 10)
	pool.Start()
	defer pool.Stop()

	// Fetch results when no tasks were submitted
	results, errs := pool.FetchResults(100 * time.Millisecond)

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
	if len(errs) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errs))
	}
}

func TestWorkerPoolV2ContextCancellation(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()

	// Submit a long-running task
	submitted := pool.Submit(func() (int, error) {
		time.Sleep(2 * time.Second)
		return 1, nil
	})
	if !submitted {
		t.Error("Failed to submit task")
	}

	// Submit another task
	submitted = pool.Submit(func() (int, error) {
		time.Sleep(2 * time.Second)
		return 2, nil
	})
	if !submitted {
		t.Error("Failed to submit task")
	}

	// Stop the pool (cancels context)
	pool.Stop()

	// Try to fetch results - should return quickly due to context cancellation
	start := time.Now()
	results, _ := pool.FetchResults(5 * time.Second)
	elapsed := time.Since(start)

	// Should complete faster than the task execution time
	if elapsed > 3*time.Second {
		t.Errorf("FetchResults took too long after context cancellation: %v", elapsed)
	}

	// We might get partial results or no results depending on timing
	if len(results) > 2 {
		t.Errorf("Expected at most 2 results, got %d", len(results))
	}
}

func TestWorkerPoolV2Counters(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 10)
	pool.Start()
	defer pool.Stop()

	// Initially all counters should be zero
	if pool.Submitted() != 0 {
		t.Errorf("Expected 0 submitted tasks, got %d", pool.Submitted())
	}
	if pool.Running() != 0 {
		t.Errorf("Expected 0 running tasks, got %d", pool.Running())
	}
	if pool.Finished() != 0 {
		t.Errorf("Expected 0 finished tasks, got %d", pool.Finished())
	}

	// Submit some tasks
	taskCount := 5
	for i := 0; i < taskCount; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 1, nil
		})
	}

	// Check submitted count
	if pool.Submitted() != taskCount {
		t.Errorf("Expected %d submitted tasks, got %d", taskCount, pool.Submitted())
	}

	// Wait a bit for tasks to start running
	time.Sleep(20 * time.Millisecond)

	// Some tasks should be running
	running := pool.Running()
	if running < 0 || running > taskCount {
		t.Errorf("Expected running tasks between 0 and %d, got %d", taskCount, running)
	}

	// Fetch all results
	results, _ := pool.FetchResults(time.Second)
	if len(results) != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, len(results))
	}

	// After fetching, submitted should be 0
	if pool.Submitted() != 0 {
		t.Errorf("Expected 0 submitted tasks after fetch, got %d", pool.Submitted())
	}

	// Running should be 0
	if pool.Running() != 0 {
		t.Errorf("Expected 0 running tasks after fetch, got %d", pool.Running())
	}

	// Finished should be 0 (decremented when fetched)
	if pool.Finished() != 0 {
		t.Errorf("Expected 0 finished tasks after fetch, got %d", pool.Finished())
	}
}

func TestWorkerPoolV2FetchAllResults(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 20)
	pool.Start()
	defer pool.Stop()

	// Submit initial batch of tasks
	for i := 0; i < 5; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 1, nil
		})
	}

	// Start fetching all results in a goroutine
	resultsCh := make(chan []int)
	go func() {
		results, _ := pool.FetchAllResults(2 * time.Second)
		resultsCh <- results
	}()

	// Submit more tasks after a delay
	time.Sleep(50 * time.Millisecond)
	for i := 0; i < 3; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 2, nil
		})
	}

	// Wait for all results
	results := <-resultsCh

	// Should get all 8 results
	if len(results) != 8 {
		t.Errorf("Expected 8 results, got %d", len(results))
	}

	// After fetching all, submitted should be 0
	if pool.Submitted() != 0 {
		t.Errorf("Expected 0 submitted tasks after fetch all, got %d", pool.Submitted())
	}
}

func TestWorkerPoolV2FetchAllResultsEmpty(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()
	defer pool.Stop()

	// Fetch when no tasks submitted
	results, errs := pool.FetchAllResults(100 * time.Millisecond)

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
	if len(errs) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errs))
	}
}

func TestWorkerPoolV2FetchAllResultsTimeout(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()
	defer pool.Stop()

	// Submit slow tasks
	for i := 0; i < 5; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(200 * time.Millisecond)
			return 1, nil
		})
	}

	// Fetch with short timeout
	start := time.Now()
	results, _ := pool.FetchAllResults(150 * time.Millisecond)
	elapsed := time.Since(start)

	// Should timeout and not get all results
	if len(results) >= 5 {
		t.Errorf("Expected fewer than 5 results due to timeout, got %d", len(results))
	}

	// Should complete around the timeout duration (with some buffer for slow systems)
	if elapsed > 400*time.Millisecond {
		t.Errorf("FetchAllResults took too long: %v", elapsed)
	}
}

func TestWorkerPoolV2FetchAllResultsWithErrors(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 10)
	pool.Start()
	defer pool.Stop()

	expectedErr := errors.New("task error")

	// Submit mix of success and error tasks
	pool.Submit(func() (int, error) {
		return 1, nil
	})
	pool.Submit(func() (int, error) {
		return 0, expectedErr
	})
	pool.Submit(func() (int, error) {
		return 2, nil
	})

	results, errs := pool.FetchAllResults(time.Second)

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
	if len(errs) != 3 {
		t.Errorf("Expected 3 error slots, got %d", len(errs))
	}

	// Check that we have the expected error
	foundError := false
	for _, err := range errs {
		if err == expectedErr {
			foundError = true
			break
		}
	}
	if !foundError {
		t.Error("Expected error not found in results")
	}
}

func TestWorkerPoolV2FetchAllResultsContinuousSubmit(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](5, 50)
	pool.Start()
	defer pool.Stop()

	// Start continuous submission in background
	stopSubmit := make(chan struct{})
	submittedCount := atomic.Int32{}
	go func() {
		for i := 0; i < 20; i++ {
			select {
			case <-stopSubmit:
				return
			default:
				pool.Submit(func() (int, error) {
					time.Sleep(10 * time.Millisecond)
					return 1, nil
				})
				submittedCount.Add(1)
				time.Sleep(15 * time.Millisecond)
			}
		}
	}()

	// Wait a bit for some tasks to be submitted
	time.Sleep(100 * time.Millisecond)

	// Stop submitting
	close(stopSubmit)

	// Fetch all results
	results, _ := pool.FetchAllResults(2 * time.Second)

	// Should get all submitted tasks
	expectedCount := int(submittedCount.Load())
	if len(results) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(results))
	}
}

func TestWorkerPoolV2CountersPrecision(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()
	defer pool.Stop()

	// Submit tasks with tracking
	taskCount := 10
	for i := 0; i < taskCount; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(20 * time.Millisecond)
			return 1, nil
		})
	}

	// Verify submitted count
	submitted := pool.Submitted()
	if submitted != taskCount {
		t.Errorf("Expected %d submitted, got %d", taskCount, submitted)
	}

	// Wait for some tasks to complete
	time.Sleep(100 * time.Millisecond)

	// Submitted should still be taskCount (not yet fetched)
	currentSubmitted := pool.Submitted()
	if currentSubmitted != taskCount {
		t.Errorf("Expected %d submitted before fetch, got %d", taskCount, currentSubmitted)
	}

	// Some tasks should be finished
	finished := pool.Finished()
	if finished == 0 {
		t.Error("Expected some finished tasks")
	}

	// Running should be 0 or more
	running := pool.Running()
	if running < 0 {
		t.Errorf("Expected non-negative running count, got %d", running)
	}

	// Fetch all results
	results, _ := pool.FetchAllResults(time.Second)
	if len(results) != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, len(results))
	}

	// All counters should be zero after fetching
	if pool.Submitted() != 0 {
		t.Errorf("Expected 0 submitted after fetch, got %d", pool.Submitted())
	}
	if pool.Running() != 0 {
		t.Errorf("Expected 0 running after fetch, got %d", pool.Running())
	}
	if pool.Finished() != 0 {
		t.Errorf("Expected 0 finished after fetch, got %d", pool.Finished())
	}
}

func TestWorkerPoolV2RunningCounter(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](2, 10)
	pool.Start()
	defer pool.Stop()

	// Submit tasks that will block
	blockCh := make(chan struct{})
	for i := 0; i < 3; i++ {
		pool.Submit(func() (int, error) {
			<-blockCh
			return 1, nil
		})
	}

	// Wait for workers to pick up tasks
	time.Sleep(50 * time.Millisecond)

	// Check running count (should be 2, limited by worker count)
	running := pool.Running()
	if running != 2 {
		t.Errorf("Expected 2 running tasks, got %d", running)
	}

	// Release one task
	blockCh <- struct{}{}
	time.Sleep(20 * time.Millisecond)

	// Third task should now be running
	running = pool.Running()
	if running != 2 {
		t.Errorf("Expected 2 running tasks after releasing one, got %d", running)
	}

	// Release all
	close(blockCh)
	time.Sleep(50 * time.Millisecond)
}

func TestWorkerPoolV2FinishedCounter(t *testing.T) {
	pool := abstract.NewWorkerPoolV2[int](3, 10)
	pool.Start()
	defer pool.Stop()

	taskCount := 5
	for i := 0; i < taskCount; i++ {
		pool.Submit(func() (int, error) {
			time.Sleep(20 * time.Millisecond)
			return 1, nil
		})
	}

	// Wait for all tasks to finish
	time.Sleep(150 * time.Millisecond)

	// All should be finished
	finished := pool.Finished()
	if finished != taskCount {
		t.Errorf("Expected %d finished tasks, got %d", taskCount, finished)
	}

	// Submitted should still be taskCount (not yet fetched)
	if pool.Submitted() != taskCount {
		t.Errorf("Expected %d submitted tasks, got %d", taskCount, pool.Submitted())
	}

	// Fetch results
	results, _ := pool.FetchResults(time.Second)
	if len(results) != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, len(results))
	}

	// After fetch, finished should be 0
	if pool.Finished() != 0 {
		t.Errorf("Expected 0 finished after fetch, got %d", pool.Finished())
	}
}

func TestWorkerPoolV2SubmitBlockingWithStop(t *testing.T) {
	// Test Submit without timeout blocking when queue is full and pool is stopped
	pool := abstract.NewWorkerPoolV2[int](1, 1)
	pool.Start()

	// Fill the queue and worker
	pool.Submit(func() (int, error) {
		time.Sleep(500 * time.Millisecond)
		return 1, nil
	})
	pool.Submit(func() (int, error) {
		return 2, nil
	})

	// Try to submit without timeout in a goroutine (it will block)
	submitDone := make(chan bool)
	go func() {
		result := pool.Submit(func() (int, error) {
			return 3, nil
		})
		submitDone <- result
	}()

	// Wait a bit to ensure submit is blocking
	time.Sleep(50 * time.Millisecond)

	// Stop the pool
	pool.Stop()

	// The submit should unblock and return false
	select {
	case result := <-submitDone:
		if result {
			t.Error("Expected submit to fail after stop")
		}
	case <-time.After(time.Second):
		t.Error("Submit did not unblock after stop")
	}
}

func TestWorkerPoolV2WorkerContextCancellation(t *testing.T) {
	// Test that worker handles context cancellation while trying to send result
	pool := abstract.NewWorkerPoolV2[int](2, 2)
	pool.Start()

	// Submit tasks that will complete but results channel might be full
	for i := 0; i < 4; i++ {
		submitted := pool.Submit(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 1, nil
		}, 100*time.Millisecond)
		if !submitted {
			// Queue is full, which is fine for this test
			break
		}
	}

	// Stop the pool while tasks are being processed
	time.Sleep(20 * time.Millisecond)
	pool.Stop()

	// Wait a bit for cleanup
	time.Sleep(100 * time.Millisecond)

	// Pool should be stopped
	if !pool.IsStopped() {
		t.Error("Pool should be stopped")
	}
}

func TestWorkerPoolV2GenericTypes(t *testing.T) {
	// Test with different types
	t.Run("String type", func(t *testing.T) {
		pool := abstract.NewWorkerPoolV2[string](2, 5)
		pool.Start()
		defer pool.Stop()

		pool.Submit(func() (string, error) {
			return "hello", nil
		})

		results, _ := pool.FetchResults(time.Second)
		if len(results) != 1 || results[0] != "hello" {
			t.Error("String type handling failed")
		}
	})

	t.Run("Struct type", func(t *testing.T) {
		type Result struct {
			ID   int
			Name string
		}

		pool := abstract.NewWorkerPoolV2[Result](2, 5)
		pool.Start()
		defer pool.Stop()

		pool.Submit(func() (Result, error) {
			return Result{ID: 1, Name: "test"}, nil
		})

		results, _ := pool.FetchResults(time.Second)
		if len(results) != 1 || results[0].ID != 1 || results[0].Name != "test" {
			t.Error("Struct type handling failed")
		}
	})

	t.Run("Pointer type", func(t *testing.T) {
		pool := abstract.NewWorkerPoolV2[*int](2, 5)
		pool.Start()
		defer pool.Stop()

		val := 42
		pool.Submit(func() (*int, error) {
			return &val, nil
		})

		results, _ := pool.FetchResults(time.Second)
		if len(results) != 1 || results[0] == nil || *results[0] != 42 {
			t.Error("Pointer type handling failed")
		}
	})
}
