package abstract

import (
	"fmt"
	"time"
)

// Timer provides precise timing measurements with support for pausing, lap timing,
// and deadline management. It's useful for performance monitoring, benchmarking,
// and creating time-based operations.
//
// Features:
//   - Precise elapsed time calculation
//   - Pause and resume functionality
//   - Lap timing for interval measurements
//   - Deadline tracking with expiration checks
//   - Multiple time unit conversions
//   - Human-readable formatting
//
// Example usage:
//
//	timer := StartTimer()
//
//	// Do some work...
//	time.Sleep(100 * time.Millisecond)
//	lap1 := timer.Lap() // Record first lap
//
//	// More work...
//	time.Sleep(200 * time.Millisecond)
//	lap2 := timer.Lap() // Record second lap
//
//	fmt.Printf("Total elapsed: %v\n", timer.ElapsedTime())
//	fmt.Printf("Formatted: %s\n", timer.FormatShort())
type Timer struct {
	start              time.Time
	paused             bool
	pausedAt           time.Time
	laps               []time.Time
	totalPauseDuration time.Duration
	deadline           time.Time
	hasDeadline        bool
}

// StartTimer creates and starts a new Timer at the current moment.
// The timer begins tracking elapsed time immediately upon creation.
//
// Returns:
//   - A new Timer instance started at the current time
//
// Example usage:
//
//	timer := StartTimer()
//
//	// Perform operations...
//	processData()
//
//	fmt.Printf("Processing took: %v\n", timer.ElapsedTime())
func StartTimer() Timer {
	return Timer{
		start: time.Now(),
		laps:  make([]time.Time, 0),
	}
}

// NewTimer creates a new Timer with the specified start time.
// This is useful for creating timers with a specific starting point.
//
// Parameters:
//   - start: The time.Time when the timer should start
//
// Returns:
//   - A new Timer instance with the specified start time
//
// Example usage:
//
//	startTime := time.Now()
//	timer := NewTimer(startTime)
//	fmt.Printf("Timer started at: %v\n", startTime)
func NewTimer(start time.Time) Timer {
	return Timer{
		start: start,
		laps:  make([]time.Time, 0),
	}
}

// String returns a human-readable string representation of the elapsed time.
// This is a convenience method that calls FormatShort() for easy output.
//
// Returns:
//   - Human-readable formatted time string
//
// Example usage:
//
//	timer := StartTimer()
//	time.Sleep(100 * time.Millisecond)
//	fmt.Printf("Elapsed: %s\n", timer.String()) // "100ms"
func (t Timer) String() string {
	return t.FormatShort()
}

// Time returns the time when the timer was started.
// This is useful for absolute time references and calculations.
//
// Returns:
//   - The time.Time when the timer was started
//
// Example usage:
//
//	timer := StartTimer()
//	startTime := timer.Time()
//	fmt.Printf("Timer started at: %v\n", startTime.Format(time.RFC3339))
func (t Timer) Time() time.Time {
	return t.start
}

// ElapsedTime returns the total time that has elapsed since the timer started,
// excluding any time spent in a paused state.
//
// Returns:
//   - Duration representing the elapsed time
//
// Example usage:
//
//	timer := StartTimer()
//	time.Sleep(100 * time.Millisecond)
//
//	elapsed := timer.ElapsedTime()
//	fmt.Printf("Elapsed: %v\n", elapsed) // ~100ms
func (t Timer) ElapsedTime() time.Duration {
	if t.paused {
		return t.pausedAt.Sub(t.start) - t.totalPauseDuration
	}
	return time.Since(t.start) - t.totalPauseDuration
}

// ElapsedSeconds returns the elapsed time as a floating-point number of seconds.
// This is convenient for calculations and when precise fractional seconds are needed.
//
// Returns:
//   - Elapsed time in seconds as a float64
//
// Example usage:
//
//	timer := StartTimer()
//	time.Sleep(1500 * time.Millisecond)
//	seconds := timer.ElapsedSeconds() // ~1.5
//	fmt.Printf("Elapsed: %.2f seconds\n", seconds)
func (t Timer) ElapsedSeconds() float64 {
	return t.ElapsedTime().Seconds()
}

// ElapsedMinutes returns the elapsed time as a floating-point number of minutes.
// Useful for longer-running operations and user-friendly time displays.
//
// Returns:
//   - Elapsed time in minutes as a float64
//
// Example usage:
//
//	timer := StartTimer()
//	// ... long-running operation ...
//	minutes := timer.ElapsedMinutes()
//	fmt.Printf("Operation took %.1f minutes\n", minutes)
func (t Timer) ElapsedMinutes() float64 {
	return t.ElapsedTime().Minutes()
}

// ElapsedHours returns the elapsed time as a floating-point number of hours.
// Suitable for very long-running operations and daily reporting.
//
// Returns:
//   - Elapsed time in hours as a float64
//
// Example usage:
//
//	timer := StartTimer()
//	// ... very long operation ...
//	hours := timer.ElapsedHours()
//	fmt.Printf("Process ran for %.2f hours\n", hours)
func (t Timer) ElapsedHours() float64 {
	return t.ElapsedTime().Hours()
}

// ElapsedMilliseconds returns the elapsed time in milliseconds.
// Useful for performance measurements and when millisecond precision is needed.
//
// Returns:
//   - Elapsed time in milliseconds as an int64
//
// Example usage:
//
//	timer := StartTimer()
//	processRequest()
//	ms := timer.ElapsedMilliseconds()
//	fmt.Printf("Request processed in %d ms\n", ms)
func (t Timer) ElapsedMilliseconds() int64 {
	return t.ElapsedTime().Milliseconds()
}

// ElapsedMicroseconds returns the elapsed time in microseconds.
// Useful for high-precision timing measurements and performance analysis.
//
// Returns:
//   - Elapsed time in microseconds as an int64
//
// Example usage:
//
//	timer := StartTimer()
//	performCriticalOperation()
//	us := timer.ElapsedMicroseconds()
//	fmt.Printf("Critical operation took %d μs\n", us)
func (t Timer) ElapsedMicroseconds() int64 {
	return t.ElapsedTime().Microseconds()
}

// ElapsedNanoseconds returns the elapsed time in nanoseconds.
// Provides the highest precision timing available for ultra-precise measurements.
//
// Returns:
//   - Elapsed time in nanoseconds as an int64
//
// Example usage:
//
//	timer := StartTimer()
//	quickOperation()
//	ns := timer.ElapsedNanoseconds()
//	fmt.Printf("Quick operation: %d ns\n", ns)
func (t Timer) ElapsedNanoseconds() int64 {
	return t.ElapsedTime().Nanoseconds()
}

// Reset resets the timer to the current time, clearing all recorded data.
// This includes laps, pause history, and deadline information.
//
// Example usage:
//
//	timer := StartTimer()
//	// ... some operations ...
//	timer.Reset() // Start timing fresh
//	// ... new operations to time ...
func (t *Timer) Reset() {
	t.start = time.Now()
	t.paused = false
	t.pausedAt = time.Time{}
	t.laps = make([]time.Time, 0)
	t.totalPauseDuration = 0
}

// Lap records the current time as a lap point and returns the duration
// since the last lap (or start time for the first lap).
//
// Returns:
//   - Duration since the previous lap or start time
//
// Example usage:
//
//	timer := StartTimer()
//
//	// Phase 1
//	doPhase1()
//	phase1Duration := timer.Lap()
//
//	// Phase 2
//	doPhase2()
//	phase2Duration := timer.Lap()
//
//	fmt.Printf("Phase 1: %v, Phase 2: %v\n", phase1Duration, phase2Duration)
func (t *Timer) Lap() time.Duration {
	now := time.Now()
	lapTime := now
	t.laps = append(t.laps, lapTime)

	if len(t.laps) == 1 {
		return lapTime.Sub(t.start)
	}
	return lapTime.Sub(t.laps[len(t.laps)-2])
}

// Laps returns a copy of all recorded lap times.
// The returned slice can be safely modified without affecting the timer.
//
// Returns:
//   - A slice of time.Time values representing when each lap was recorded
//
// Example usage:
//
//	timer := StartTimer()
//	timer.Lap() // Lap 1
//	timer.Lap() // Lap 2
//
//	laps := timer.Laps()
//	for i, lap := range laps {
//		fmt.Printf("Lap %d recorded at: %v\n", i+1, lap)
//	}
func (t Timer) Laps() []time.Time {
	// Return a copy to prevent modification
	result := make([]time.Time, len(t.laps))
	copy(result, t.laps)
	return result
}

// LapDurations returns the durations between consecutive laps.
// The first duration is measured from the start time to the first lap.
//
// Returns:
//   - A slice of durations representing the time for each lap interval
//
// Example usage:
//
//	timer := StartTimer()
//	// ... work phase 1 ...
//	timer.Lap()
//	// ... work phase 2 ...
//	timer.Lap()
//
//	durations := timer.LapDurations()
//	for i, duration := range durations {
//		fmt.Printf("Lap %d duration: %v\n", i+1, duration)
//	}
func (t Timer) LapDurations() []time.Duration {
	if len(t.laps) == 0 {
		return nil
	}

	durations := make([]time.Duration, len(t.laps))
	durations[0] = t.laps[0].Sub(t.start)

	for i := 1; i < len(t.laps); i++ {
		durations[i] = t.laps[i].Sub(t.laps[i-1])
	}

	return durations
}

// Format returns the elapsed time formatted according to a custom layout.
// The layout uses Go's standard time formatting with placeholders for
// hours, minutes, seconds, and milliseconds.
//
// Parameters:
//   - layout: Format string with placeholders for time components
//
// Returns:
//   - Formatted string representation of elapsed time
//
// Example usage:
//
//	timer := StartTimer()
//	time.Sleep(1*time.Hour + 23*time.Minute + 45*time.Second + 678*time.Millisecond)
//
//	formatted := timer.Format("%02d:%02d:%02d.%03d") // "01:23:45.678"
//	fmt.Printf("Elapsed: %s\n", formatted)
func (t Timer) Format(layout string) string {
	elapsed := t.ElapsedTime()
	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60
	milliseconds := int(elapsed.Milliseconds()) % 1000

	return fmt.Sprintf(layout, hours, minutes, seconds, milliseconds)
}

// FormatShort returns a human-readable string representation of elapsed time.
// The format automatically adjusts based on the duration magnitude for optimal readability.
//
// Format examples:
//   - Less than 1 second: "123ms"
//   - Less than 1 minute: "12.34s"
//   - Less than 1 hour: "12m34s"
//   - 1 hour or more: "1h23m45s"
//
// Returns:
//   - Human-readable formatted time string
//
// Example usage:
//
//	timer := StartTimer()
//	time.Sleep(2*time.Minute + 30*time.Second)
//	fmt.Printf("Elapsed: %s\n", timer.FormatShort()) // "2m30s"
func (t Timer) FormatShort() string {
	d := t.ElapsedTime()

	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	} else if d < time.Hour {
		s := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), s)
	}

	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%dh%dm%ds", h, m, s)
}

// HasElapsed checks if a specified duration has elapsed since the timer started.
// This is useful for implementing timeouts and periodic checks.
//
// Parameters:
//   - duration: The duration to check against
//
// Returns:
//   - true if the specified duration has elapsed, false otherwise
//
// Example usage:
//
//	timer := StartTimer()
//
//	for !timer.HasElapsed(30 * time.Second) {
//		// Do work until 30 seconds have elapsed
//		doWork()
//		time.Sleep(100 * time.Millisecond)
//	}
//	fmt.Println("30 seconds have passed!")
func (t Timer) HasElapsed(duration time.Duration) bool {
	return t.ElapsedTime() >= duration
}

// Pause pauses the timer, stopping the accumulation of elapsed time.
// Subsequent calls to elapsed time methods will not include time after the pause
// until Resume() is called.
//
// Returns:
//   - true if the timer was successfully paused, false if already paused
//
// Example usage:
//
//	timer := StartTimer()
//	// ... do work ...
//
//	if timer.Pause() {
//		fmt.Println("Timer paused")
//		// ... timer is not counting during this period ...
//		time.Sleep(5 * time.Second) // This won't count toward elapsed time
//
//		timer.Resume()
//		fmt.Println("Timer resumed")
//	}
func (t *Timer) Pause() bool {
	if t.paused {
		return false
	}
	t.paused = true
	t.pausedAt = time.Now()
	return true
}

// Resume resumes the timer if it was paused, continuing the accumulation
// of elapsed time from where it left off.
//
// Returns:
//   - true if the timer was successfully resumed, false if not paused
//
// Example usage:
//
//	timer := StartTimer()
//	timer.Pause()
//	// ... do non-timed work ...
//
//	if timer.Resume() {
//		fmt.Println("Timer resumed, continuing to track time")
//		// ... timer continues counting elapsed time ...
//	}
func (t *Timer) Resume() bool {
	if !t.paused {
		return false
	}

	pauseDuration := time.Since(t.pausedAt)
	t.totalPauseDuration += pauseDuration
	t.paused = false
	return true
}

// IsPaused returns whether the timer is currently in a paused state.
//
// Returns:
//   - true if the timer is paused, false if actively running
//
// Example usage:
//
//	timer := StartTimer()
//	timer.Pause()
//
//	if timer.IsPaused() {
//		fmt.Println("Timer is currently paused")
//		timer.Resume()
//	}
func (t Timer) IsPaused() bool {
	return t.paused
}

// Deadline creates a new timer with a deadline set to the specified duration from now.
// This is useful for creating timers with built-in expiration functionality.
//
// Parameters:
//   - duration: Time duration until the deadline
//
// Returns:
//   - A new Timer with the deadline set
//
// Example usage:
//
//	// Create a timer that expires in 5 minutes
//	timer := Deadline(5 * time.Minute)
//
//	for !timer.IsExpired() {
//		// Do work until deadline is reached
//		doWork()
//		time.Sleep(1 * time.Second)
//
//		remaining := timer.TimeRemaining()
//		fmt.Printf("Time remaining: %v\n", remaining)
//	}
//	fmt.Println("Deadline reached!")
func Deadline(duration time.Duration) Timer {
	t := StartTimer()
	t.deadline = time.Now().Add(duration)
	t.hasDeadline = true
	return t
}

// SetDeadline sets an absolute deadline time for the timer.
// This allows for precise deadline management with specific target times.
//
// Parameters:
//   - deadline: The absolute time when the deadline occurs
//
// Example usage:
//
//	timer := StartTimer()
//
//	// Set deadline to midnight tomorrow
//	tomorrow := time.Now().Add(24 * time.Hour)
//	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
//	timer.SetDeadline(midnight)
//
//	fmt.Printf("Deadline set for: %v\n", midnight)
func (t *Timer) SetDeadline(deadline time.Time) {
	t.deadline = deadline
	t.hasDeadline = true
}

// SetDeadlineDuration sets a deadline relative to the current time.
// This is equivalent to calling SetDeadline(time.Now().Add(duration)).
//
// Parameters:
//   - duration: Duration from now until the deadline
//
// Example usage:
//
//	timer := StartTimer()
//	timer.SetDeadlineDuration(10 * time.Minute) // Deadline in 10 minutes
//
//	// Check deadline status periodically
//	for !timer.IsExpired() {
//		processTask()
//		fmt.Printf("Time remaining: %v\n", timer.TimeRemaining())
//	}
func (t *Timer) SetDeadlineDuration(duration time.Duration) {
	t.deadline = time.Now().Add(duration)
	t.hasDeadline = true
}

// TimeRemaining returns the time remaining until the deadline.
// If no deadline is set or the deadline has passed, returns zero duration.
//
// Returns:
//   - Duration remaining until deadline, or 0 if expired/no deadline
//
// Example usage:
//
//	timer := Deadline(30 * time.Second)
//
//	for timer.TimeRemaining() > 0 {
//		remaining := timer.TimeRemaining()
//		fmt.Printf("Countdown: %v\n", remaining)
//		time.Sleep(1 * time.Second)
//	}
//	fmt.Println("Time's up!")
func (t Timer) TimeRemaining() time.Duration {
	if !t.hasDeadline {
		return 0
	}

	remaining := time.Until(t.deadline)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsExpired returns true if the deadline has passed.
// If no deadline is set, always returns false.
//
// Returns:
//   - true if deadline has been reached, false otherwise
//
// Example usage:
//
//	timer := Deadline(1 * time.Minute)
//
//	// Process items until deadline
//	for !timer.IsExpired() {
//		if processNextItem() {
//			fmt.Printf("Processed item, %v remaining\n", timer.TimeRemaining())
//		} else {
//			break // No more items
//		}
//	}
//
//	if timer.IsExpired() {
//		fmt.Println("Deadline reached, stopping processing")
//	}
func (t Timer) IsExpired() bool {
	if !t.hasDeadline {
		return false
	}
	return time.Now().After(t.deadline)
}
