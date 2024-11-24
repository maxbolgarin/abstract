package abstract

import (
	"context"
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
