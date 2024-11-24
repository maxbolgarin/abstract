package abstract

import (
	"context"
	"errors"
	"time"

	"github.com/maxbolgarin/lang"
)

var ErrTimeout = errors.New("timeout")

// Future is used for running a function in a separate goroutine with returning the result.
//
// How to use:
//
//	f1 := abstract.NewFuture(ctx, slog.Default(), func(context.Context) (string, error) {
//		// TODO: some code
//		return "some result", nil
//	})
//
//	result, err := f1.Get(ctx)
type Future[T any] struct {
	value T
	err   error
	done  chan struct{}
}

// NewFuture returns a new started future, it creates a goroutine that will run the passed function and remember
// it's result and error.
func NewFuture[T any](ctx context.Context, l lang.Logger, foo func(ctx context.Context) (T, error)) *Future[T] {
	future := &Future[T]{
		done: make(chan struct{}),
	}

	go func() {
		defer close(future.done)
		defer lang.RecoverWithErrAndStack(l, &future.err)
		future.value, future.err = foo(ctx)
	}()

	return future
}

// Get will wait for the result of the underlying future or returns without it if the context is canceled.
func (f *Future[T]) Get(ctx context.Context) (T, error) {
	// Firstly try to get result without checking the context (need for WaiterSet).
	select {
	case <-f.done:
		return f.value, f.err
	default:
	}

	select {
	case <-ctx.Done():
		var empty T
		return empty, ctx.Err()
	case <-f.done:
		return f.value, f.err
	}
}

// GetWithTimeout will wait for the result of the underlying future or returns without it if the context is canceled,
// it also wait for the result for the provided timeout.
func (f *Future[T]) GetWithTimeout(ctx context.Context, timeout time.Duration) (T, error) {
	// Firstly try to get result without checking the context and timeout (need for WaiterSet).
	select {
	case <-f.done:
		return f.value, f.err
	default:
	}

	select {
	case <-ctx.Done():
		var empty T
		return empty, ctx.Err()
	case <-time.After(timeout):
		var empty T
		return empty, ErrTimeout
	case <-f.done:
		return f.value, f.err
	}
}

// Waiter is used for running a function in a separate goroutine with returning the error.
//
// How to use:
//
//	w := abstract.NewWaiter(ctx, slog.Default(), func(context.Context) error {
//		// TODO: some code
//		return nil
//	})
//
//	err := w.Await(ctx)
type Waiter struct {
	f *Future[struct{}]
}

// NewFutureVoid returns a new started void future, it creates a goroutine that will run the passed function
// and remember it's error.
func NewWaiter(ctx context.Context, l lang.Logger, foo func(context.Context) error) *Waiter {
	return &Waiter{
		f: NewFuture(ctx, l, func(ctx context.Context) (struct{}, error) {
			return struct{}{}, foo(ctx)
		}),
	}
}

// Await will wait for the result of the underlying future or returns without it if the context is canceled.
func (f *Waiter) Await(ctx context.Context) error {
	_, err := f.f.Get(ctx)
	return err
}

// AwaitWithTimeout will wait for the result of the underlying future or returns without it if the context is canceled,
// it also wait for the result for the provided timeout.
func (f *Waiter) AwaitWithTimeout(ctx context.Context, timeout time.Duration) error {
	_, err := f.f.GetWithTimeout(ctx, timeout)
	return err
}

// WaiterSet is used for running many functions, each in a separate goroutine
// with returning a combined error.
//
// How to use:
//
//	ws := abstract.NewWaiterSet(slog.Default())
//	ws.Add(ctx, func(context.Context) error {
//		// TODO: some code 1
//		return nil
//	})
//	ws.Add(ctx, func(context.Context) error {
//		// TODO: some code 2
//		return nil
//	})
//
//	err := ws.Await(ctx) // Wait for completion of all added functions
type WaiterSet struct {
	ws []*Waiter
	l  lang.Logger
}

// NewWaiterSet returns new [WaiterSet] with added [Waiter], that were started earlier.
func NewWaiterSet(l lang.Logger, ws ...*Waiter) *WaiterSet {
	return &WaiterSet{
		ws: ws,
		l:  l,
	}
}

// Add adds a new [Waiter] to the [WaiterSet] and starts it in a separate goroutine.
func (s *WaiterSet) Add(ctx context.Context, foo func(ctx context.Context) error) {
	s.ws = append(s.ws, NewWaiter(ctx, s.l, foo))
}

// Await will wait for the result of all underlying functions or returns without it if the context is canceled.
func (s *WaiterSet) Await(ctx context.Context) error {
	var errs []error
	for _, w := range s.ws {
		err := w.Await(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// AwaitWithTimeout will wait for the result of all underlying functions or returns without it if the context is canceled,
// it also wait for the result for the provided timeout.
func (s *WaiterSet) AwaitWithTimeout(ctx context.Context, timeout time.Duration) error {
	var (
		errs []error
		t    = StartTimer()
	)

	for _, w := range s.ws {
		currentTimeout := timeout - t.ElapsedTime()
		if currentTimeout < 0 {
			currentTimeout = 0
		}
		err := w.AwaitWithTimeout(ctx, currentTimeout)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
