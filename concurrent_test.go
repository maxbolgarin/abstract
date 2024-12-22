package abstract_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

// TestStartUpdater tests the StartUpdater function.
func TestStartUpdater(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count atomic.Int64
	interval := 20 * time.Millisecond
	f := func() {
		count.Add(1)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	abstract.StartUpdater(ctx, interval, nil, f)
	time.Sleep(100 * time.Millisecond)

	if count.Load() < 2 || count.Load() > 3 {
		t.Errorf("expected function to be called 2 or 3 times, got %d", count.Load())
	}
}

// TestStartUpdaterNow tests the StartUpdaterNow function.
func TestStartUpdaterNow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count atomic.Int64
	interval := 20 * time.Millisecond
	f := func() {
		count.Add(1)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	abstract.StartUpdaterNow(ctx, interval, nil, f)
	time.Sleep(100 * time.Millisecond)

	if count.Load() < 3 || count.Load() > 4 {
		t.Errorf("expected function to be called 3 or 4 times, got %d", count.Load())
	}
}

// TestStartUpdaterWithShutdown tests the StartUpdaterWithShutdown function.
func TestStartUpdaterWithShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count atomic.Int64
	var shutdownCalled atomic.Bool
	interval := 20 * time.Millisecond
	f := func() {
		count.Add(1)
	}
	shutdown := func() {
		shutdownCalled.Store(true)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	abstract.StartUpdaterWithShutdown(ctx, interval, nil, f, shutdown)
	time.Sleep(100 * time.Millisecond)

	if count.Load() < 2 || count.Load() > 3 {
		t.Errorf("expected function to be called 2 or 3 times, got %d", count.Load())
	}
	if !shutdownCalled.Load() {
		t.Errorf("expected shutdown to be called")
	}
}

// TestStartUpdaterWithShutdownChan tests the StartUpdaterWithShutdownChan function.
func TestStartUpdaterWithShutdownChan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count atomic.Int64
	interval := 20 * time.Millisecond
	ch := make(chan struct{})
	f := func() {
		count.Add(1)
	}

	go func() {
		time.Sleep(30 * time.Millisecond)
		close(ch)
	}()

	abstract.StartUpdaterWithShutdownChan(ctx, interval, nil, ch, f)
	time.Sleep(100 * time.Millisecond)

	if count.Load() < 1 || count.Load() > 2 {
		t.Errorf("expected function to be called 1 or 2 times, got %d", count.Load())
	}
}

func TestRateProcessor(t *testing.T) {
	ctx := context.Background()
	rp := abstract.NewRateProcessor(ctx, 5)

	taskCount := 5
	for i := 0; i < taskCount; i++ {
		rp.AddTask(func(context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}

	errors := rp.Wait()
	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %v", len(errors))
	}
}

func TestRateProcessorWithErrors(t *testing.T) {
	ctx := context.Background()
	rp := abstract.NewRateProcessor(ctx, 5)

	taskCount := 5
	for i := 0; i < taskCount; i++ {
		rp.AddTask(func(context.Context) error {
			return errors.New("task error")
		})
	}

	errors := rp.Wait()
	if len(errors) != taskCount {
		t.Errorf("Expected %v errors, got %v", taskCount, len(errors))
	}
}

func TestRateProcessorCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rp := abstract.NewRateProcessor(ctx, 2)

	cancel() // Отмена контекста сразу после создания

	rp.AddTask(func(context.Context) error {
		return nil
	})

	errors := rp.Wait()
	if len(errors) != 0 {
		t.Errorf("Expected no errors due to immediate cancellation, got %v", len(errors))
	}
}
