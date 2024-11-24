package abstract

import "time"

// Timer is used to calculate time intervals.
type Timer struct {
	start time.Time
}

// Start starts [Timer] at the current moment of time.
func StartTimer() Timer {
	return Timer{start: time.Now()}
}

// Time returns the startup time.
func (t Timer) Time() time.Time {
	return t.start
}

// ElapsedTime returns the time interval before [Timer] starting.
func (t Timer) ElapsedTime() time.Duration {
	return time.Since(t.start)
}
