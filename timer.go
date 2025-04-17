package abstract

import (
	"fmt"
	"time"
)

// Timer is used to calculate time intervals.
type Timer struct {
	start              time.Time
	paused             bool
	pausedAt           time.Time
	laps               []time.Time
	totalPauseDuration time.Duration
	deadline           time.Time
	hasDeadline        bool
}

// StartTimer starts a new Timer at the current moment of time.
func StartTimer() Timer {
	return Timer{
		start: time.Now(),
		laps:  make([]time.Time, 0),
	}
}

// Time returns the startup time.
func (t Timer) Time() time.Time {
	return t.start
}

// ElapsedTime returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedTime() time.Duration {
	if t.paused {
		return t.pausedAt.Sub(t.start) - t.totalPauseDuration
	}
	return time.Since(t.start) - t.totalPauseDuration
}

// ElapsedSeconds returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedSeconds() float64 {
	return t.ElapsedTime().Seconds()
}

// ElapsedMinutes returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedMinutes() float64 {
	return t.ElapsedTime().Minutes()
}

// ElapsedHours returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedHours() float64 {
	return t.ElapsedTime().Hours()
}

// ElapsedMilliseconds returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedMilliseconds() int64 {
	return t.ElapsedTime().Milliseconds()
}

// ElapsedMicroseconds returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedMicroseconds() int64 {
	return t.ElapsedTime().Microseconds()
}

// ElapsedNanoseconds returns the time interval since Timer starting,
// accounting for any paused time.
func (t Timer) ElapsedNanoseconds() int64 {
	return t.ElapsedTime().Nanoseconds()
}

// Reset resets the timer to the current time and clears all laps and pause history.
func (t *Timer) Reset() {
	t.start = time.Now()
	t.paused = false
	t.pausedAt = time.Time{}
	t.laps = make([]time.Time, 0)
	t.totalPauseDuration = 0
}

// Lap records the current time as a lap and returns the lap duration.
func (t *Timer) Lap() time.Duration {
	now := time.Now()
	lapTime := now
	t.laps = append(t.laps, lapTime)

	if len(t.laps) == 1 {
		return lapTime.Sub(t.start)
	}
	return lapTime.Sub(t.laps[len(t.laps)-2])
}

// Laps returns all recorded lap times.
func (t Timer) Laps() []time.Time {
	// Return a copy to prevent modification
	result := make([]time.Time, len(t.laps))
	copy(result, t.laps)
	return result
}

// LapDurations returns the durations between consecutive laps.
// The first duration is measured from the start time.
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

// Format returns the elapsed time formatted according to the given layout.
// Common layouts: "15:04:05", "04:05.000", etc.
func (t Timer) Format(layout string) string {
	elapsed := t.ElapsedTime()
	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60
	milliseconds := int(elapsed.Milliseconds()) % 1000

	return fmt.Sprintf(layout, hours, minutes, seconds, milliseconds)
}

// FormatShort returns a human-readable string representation of the elapsed time.
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

// HasElapsed checks if the specified duration has elapsed since the timer started.
func (t Timer) HasElapsed(duration time.Duration) bool {
	return t.ElapsedTime() >= duration
}

// Pause pauses the timer. Subsequent calls to ElapsedTime will not include time after the pause.
// Returns false if the timer is already paused.
func (t *Timer) Pause() bool {
	if t.paused {
		return false
	}
	t.paused = true
	t.pausedAt = time.Now()
	return true
}

// Resume resumes the timer if it was paused.
// Returns false if the timer was not paused.
func (t *Timer) Resume() bool {
	if !t.paused {
		return false
	}

	pauseDuration := time.Since(t.pausedAt)
	t.totalPauseDuration += pauseDuration
	t.paused = false
	return true
}

// IsPaused returns whether the timer is currently paused.
func (t Timer) IsPaused() bool {
	return t.paused
}

// Deadline creates a new timer with a deadline set to the specified duration from now.
func Deadline(duration time.Duration) Timer {
	t := StartTimer()
	t.deadline = time.Now().Add(duration)
	t.hasDeadline = true
	return t
}

// SetDeadline sets a deadline for the timer.
func (t *Timer) SetDeadline(deadline time.Time) {
	t.deadline = deadline
	t.hasDeadline = true
}

// SetDeadlineDuration sets a deadline relative to the current time.
func (t *Timer) SetDeadlineDuration(duration time.Duration) {
	t.deadline = time.Now().Add(duration)
	t.hasDeadline = true
}

// TimeRemaining returns the time remaining until the deadline.
// If no deadline is set or the deadline has passed, it returns zero duration.
func (t Timer) TimeRemaining() time.Duration {
	if !t.hasDeadline {
		return 0
	}

	remaining := t.deadline.Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsExpired returns true if the deadline has passed.
// If no deadline is set, it returns false.
func (t Timer) IsExpired() bool {
	if !t.hasDeadline {
		return false
	}
	return time.Now().After(t.deadline)
}
