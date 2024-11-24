package abstract_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

func TestFuture(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f1 := abstract.NewFuture(ctx, nil, func(context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 10, errors.New("real error")
	})

	result, err := f1.GetWithTimeout(ctx, time.Millisecond)
	if result != 0 {
		t.Errorf("expected 0 but got %d", result)
	}
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error but got %v", err)
	}

	cancel()
	result, err = f1.Get(ctx)
	if result != 0 {
		t.Errorf("expected 0 but got %d", result)
	}
	if err != ctx.Err() {
		t.Errorf("expected context canceled error but got %v", err)
	}

	result, err = f1.GetWithTimeout(ctx, 100*time.Millisecond)
	if result != 0 {
		t.Errorf("expected 0 but got %d", result)
	}
	if err != ctx.Err() {
		t.Errorf("expected context canceled error but got %v", err)
	}

	result, err = f1.Get(context.Background())
	if result != 10 {
		t.Errorf("expected 10 but got %d", result)
	}
	if err == nil || !strings.Contains(err.Error(), "real error") {
		t.Errorf("expected real error but got %v", err)
	}

	result, err = f1.GetWithTimeout(context.Background(), time.Nanosecond)
	if result != 10 {
		t.Errorf("expected 10 but got %d", result)
	}
	if err == nil || !strings.Contains(err.Error(), "real error") {
		t.Errorf("expected real error but got %v", err)
	}
}

func TestWaiter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f1 := abstract.NewWaiter(ctx, nil, func(context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return errors.New("real error")
	})

	err := f1.AwaitWithTimeout(ctx, time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error but got %v", err)
	}

	cancel()
	err = f1.Await(ctx)
	if err != ctx.Err() {
		t.Errorf("expected context canceled error but got %v", err)
	}

	err = f1.AwaitWithTimeout(ctx, 100*time.Millisecond)
	if err != ctx.Err() {
		t.Errorf("expected context canceled error but got %v", err)
	}

	err = f1.Await(context.Background())
	if err == nil || !strings.Contains(err.Error(), "real error") {
		t.Errorf("expected real error but got %v", err)
	}

	err = f1.AwaitWithTimeout(context.Background(), time.Nanosecond)
	if err == nil || !strings.Contains(err.Error(), "real error") {
		t.Errorf("expected real error but got %v", err)
	}
}

func TestWaiterSet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws := abstract.NewWaiterSet(nil)

	ws.Add(ctx, func(context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return errors.New("error1")
	})
	ws.Add(ctx, func(context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return errors.New("error2")
	})
	ws.Add(ctx, func(context.Context) error {
		time.Sleep(10 * time.Millisecond)
		panic(errors.New("error3"))
	})
	for range [10]int{} {
		ws.Add(ctx, func(context.Context) error {
			return nil
		})
	}

	err := ws.AwaitWithTimeout(ctx, time.Millisecond)

	if err == nil || !strings.Contains(err.Error(), abstract.ErrTimeout.Error()) {
		t.Errorf("expected timeout error but got %v", err)
	}

	cancel()
	err = ws.Await(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error but got %v", err)
	}

	err = ws.AwaitWithTimeout(ctx, 100*time.Millisecond)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error but got %v", err)
	}

	err = ws.AwaitWithTimeout(context.Background(), 60*time.Millisecond)
	if !strings.Contains(err.Error(), errors.New("error2").Error()) {
		t.Errorf("expected error2 but got %v", err)
	}
	if !strings.Contains(err.Error(), errors.New("error3").Error()) {
		t.Errorf("expected error3 but got %v", err)
	}
	if !strings.Contains(err.Error(), abstract.ErrTimeout.Error()) {
		t.Errorf("expected timeout error but got %v", err)
	}

	err = ws.Await(context.Background())
	if !strings.Contains(err.Error(), errors.New("error1").Error()) {
		t.Errorf("expected error1 but got %v", err)
	}
	if !strings.Contains(err.Error(), errors.New("error2").Error()) {
		t.Errorf("expected error2 but got %v", err)
	}
	if !strings.Contains(err.Error(), errors.New("error3").Error()) {
		t.Errorf("expected error3 but got %v", err)
	}
	if strings.Contains(err.Error(), abstract.ErrTimeout.Error()) {
		t.Errorf("did not expect timeout error but got %v", err)
	}
}
