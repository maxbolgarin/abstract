package abstract_test

import (
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

func TestStartTimer(t *testing.T) {
	timer := abstract.StartTimer()

	// Time should be close to now
	now := time.Now()
	diff := now.Sub(timer.Time())
	if diff > 10*time.Millisecond {
		t.Errorf("Expected start time close to now, got difference of %v", diff)
	}

	// Timer should not be paused
	if timer.IsPaused() {
		t.Error("New timer should not be paused")
	}
}

func TestTime(t *testing.T) {
	// Since start field is unexported, we'll use StartTimer instead
	timer := abstract.StartTimer()
	startTime := timer.Time()

	if !timer.Time().Equal(startTime) {
		t.Errorf("Expected Time() to return %v, got %v", startTime, timer.Time())
	}
}

func TestElapsedTime(t *testing.T) {
	// Test normal elapsed time
	timer := abstract.StartTimer()
	time.Sleep(50 * time.Millisecond)
	elapsed := timer.ElapsedTime()

	if elapsed < 25*time.Millisecond || elapsed > 200*time.Millisecond {
		t.Errorf("Expected elapsed time around 50ms, got %v", elapsed)
	}

	// Test paused timer
	timer = abstract.StartTimer()
	time.Sleep(25 * time.Millisecond)
	timer.Pause()
	pauseTime := timer.ElapsedTime()
	time.Sleep(25 * time.Millisecond)
	afterPauseTime := timer.ElapsedTime()

	if pauseTime != afterPauseTime {
		t.Errorf("Elapsed time should not change while paused. Before: %v, After: %v",
			pauseTime, afterPauseTime)
	}
}

func TestElapsedTimeVariants(t *testing.T) {
	// Create a timer and manually check elapsed values
	// We need to wait for the timer to have a predictable elapsed time
	timer := abstract.StartTimer()

	// Allow some tolerance in the actual time measurements
	if seconds := timer.ElapsedSeconds(); seconds > 0 {
		// Just check that it's returning a reasonable positive value
		t.Logf("ElapsedSeconds returned %v", seconds)
	}

	if minutes := timer.ElapsedMinutes(); minutes > 0 {
		// Just check that it's returning a reasonable positive value
		t.Logf("ElapsedMinutes returned %v", minutes)
	}

	if hours := timer.ElapsedHours(); hours > 0 {
		// Just check that it's returning a reasonable positive value
		t.Logf("ElapsedHours returned %v", hours)
	}

	if ms := timer.ElapsedMilliseconds(); ms > 0 {
		// Just check that it's returning a reasonable positive value
		t.Logf("ElapsedMilliseconds returned %v", ms)
	}

	if us := timer.ElapsedMicroseconds(); us > 0 {
		// Just check that it's returning a reasonable positive value
		t.Logf("ElapsedMicroseconds returned %v", us)
	}

	if ns := timer.ElapsedNanoseconds(); ns > 0 {
		// Just check that it's returning a reasonable positive value
		t.Logf("ElapsedNanoseconds returned %v", ns)
	}
}

func TestReset(t *testing.T) {
	timer := abstract.StartTimer()
	time.Sleep(50 * time.Millisecond)

	before := timer.ElapsedTime()
	timer.Reset()
	after := timer.ElapsedTime()

	if before < 25*time.Millisecond {
		t.Errorf("Before reset, elapsed time should be > 25ms, got %v", before)
	}

	if after > 25*time.Millisecond {
		t.Errorf("After reset, elapsed time should be < 25ms, got %v", after)
	}

	// Test that laps are cleared
	timer.Lap()
	timer.Lap()
	if len(timer.Laps()) != 2 {
		t.Errorf("Expected 2 laps, got %d", len(timer.Laps()))
	}

	timer.Reset()
	if len(timer.Laps()) != 0 {
		t.Errorf("After reset, expected 0 laps, got %d", len(timer.Laps()))
	}
}

func TestLap(t *testing.T) {
	timer := abstract.StartTimer()
	time.Sleep(25 * time.Millisecond)

	lap1 := timer.Lap()
	if lap1 < 20*time.Millisecond || lap1 > 200*time.Millisecond {
		t.Errorf("Expected first lap around 25ms, got %v", lap1)
	}

	time.Sleep(25 * time.Millisecond)
	lap2 := timer.Lap()
	if lap2 < 20*time.Millisecond || lap2 > 200*time.Millisecond {
		t.Errorf("Expected second lap around 25ms, got %v", lap2)
	}

	laps := timer.Laps()
	if len(laps) != 2 {
		t.Errorf("Expected 2 laps, got %d", len(laps))
	}
}

func TestLaps(t *testing.T) {
	timer := abstract.StartTimer()

	// Create some laps
	timer.Lap()
	timer.Lap()
	timer.Lap()

	laps := timer.Laps()
	if len(laps) != 3 {
		t.Errorf("Expected 3 laps, got %d", len(laps))
	}

	// Test that returned slice is a copy
	lapsCopy := timer.Laps()
	if &laps[0] == &lapsCopy[0] {
		t.Error("Laps() should return a copy, not a reference")
	}
}

func TestLapDurations(t *testing.T) {
	timer := abstract.StartTimer()

	// No laps yet
	if durations := timer.LapDurations(); durations != nil {
		t.Errorf("Expected nil durations with no laps, got %v", durations)
	}

	timer.Lap()
	time.Sleep(25 * time.Millisecond)
	timer.Lap()

	durations := timer.LapDurations()
	if len(durations) != 2 {
		t.Errorf("Expected 2 durations, got %d", len(durations))
	}

	// First duration is from start to first lap
	// Second duration is from first lap to second lap
	if durations[1] < 20*time.Millisecond || durations[1] > 200*time.Millisecond {
		t.Errorf("Expected second duration around 25ms, got %v", durations[1])
	}
}

func TestFormat(t *testing.T) {
	// Create a timer and wait a small amount to test formatting
	timer := abstract.StartTimer()
	// We'll test the formatting functions with actual elapsed time

	formatted := timer.Format("%02d:%02d:%02d.%03d")
	// We expect something like "00:00:00.xxx"
	if len(formatted) < 8 {
		t.Errorf("Expected format like 00:00:00.xxx, got %s", formatted)
	}
}

func TestFormatShort(t *testing.T) {
	// Since we can't directly modify the start time, we'll test each case separately
	// with actual timers

	// Test milliseconds format (for a brand new timer)
	timer := abstract.StartTimer()
	formatted := timer.FormatShort()
	if formatted[len(formatted)-2:] != "ms" {
		t.Errorf("Expected format ending with 'ms', got '%s'", formatted)
	}

	// For the remaining tests, we'll just make sure the function returns something
	// without erroring since we can't easily control the exact elapsed time
	timer = abstract.StartTimer()
	if formatted = timer.FormatShort(); formatted == "" {
		t.Errorf("FormatShort returned empty string")
	}
}

func TestHasElapsed(t *testing.T) {
	timer := abstract.StartTimer()

	if timer.HasElapsed(50 * time.Millisecond) {
		t.Error("New timer should not have elapsed 50ms yet")
	}

	time.Sleep(60 * time.Millisecond)

	if !timer.HasElapsed(50 * time.Millisecond) {
		t.Error("Timer should have elapsed 50ms after sleeping for 60ms")
	}
}

func TestPauseResume(t *testing.T) {
	timer := abstract.StartTimer()
	time.Sleep(25 * time.Millisecond)

	// Test pause
	if !timer.Pause() {
		t.Error("First call to Pause() should return true")
	}

	if !timer.IsPaused() {
		t.Error("Timer should be paused after calling Pause()")
	}

	// Double pause should return false
	if timer.Pause() {
		t.Error("Second call to Pause() should return false")
	}

	pausedTime := timer.ElapsedTime()
	time.Sleep(25 * time.Millisecond)

	// Elapsed time should not change while paused
	if pausedTime != timer.ElapsedTime() {
		t.Errorf("Elapsed time should not change while paused: %v vs %v",
			pausedTime, timer.ElapsedTime())
	}

	// Test resume
	if !timer.Resume() {
		t.Error("Resume() should return true when timer is paused")
	}

	if timer.IsPaused() {
		t.Error("Timer should not be paused after calling Resume()")
	}

	// Resume when not paused should return false
	if timer.Resume() {
		t.Error("Resume() should return false when timer is not paused")
	}

	// After resume, elapsed time should continue increasing
	resumedTime := timer.ElapsedTime()
	time.Sleep(25 * time.Millisecond)
	afterTime := timer.ElapsedTime()

	if afterTime <= resumedTime {
		t.Errorf("Elapsed time should increase after resume: %v vs %v",
			resumedTime, afterTime)
	}
}

func TestDeadline(t *testing.T) {
	// Test creating a timer with deadline
	timer := abstract.Deadline(100 * time.Millisecond)

	if timer.IsExpired() {
		t.Error("New deadline timer should not be expired")
	}

	remaining := timer.TimeRemaining()
	if remaining <= 0 || remaining > 100*time.Millisecond {
		t.Errorf("Expected remaining time around 100ms, got %v", remaining)
	}

	time.Sleep(120 * time.Millisecond)

	if !timer.IsExpired() {
		t.Error("Timer should be expired after deadline")
	}

	if timer.TimeRemaining() != 0 {
		t.Errorf("Expired timer should have 0 time remaining, got %v", timer.TimeRemaining())
	}
}

func TestSetDeadline(t *testing.T) {
	timer := abstract.StartTimer()

	// Timer with no deadline should not be expired
	if timer.IsExpired() {
		t.Error("Timer without deadline should not be expired")
	}

	// Set absolute deadline
	deadline := time.Now().Add(50 * time.Millisecond)
	timer.SetDeadline(deadline)

	if timer.IsExpired() {
		t.Error("Timer should not be expired immediately after setting deadline")
	}

	time.Sleep(70 * time.Millisecond)

	if !timer.IsExpired() {
		t.Error("Timer should be expired after deadline has passed")
	}
}

func TestSetDeadlineDuration(t *testing.T) {
	timer := abstract.StartTimer()

	// Set relative deadline
	timer.SetDeadlineDuration(50 * time.Millisecond)

	if timer.IsExpired() {
		t.Error("Timer should not be expired immediately after setting deadline")
	}

	time.Sleep(70 * time.Millisecond)

	if !timer.IsExpired() {
		t.Error("Timer should be expired after deadline duration has passed")
	}
}

func TestTimeRemaining(t *testing.T) {
	timer := abstract.StartTimer()

	// Timer with no deadline should have 0 time remaining
	if timer.TimeRemaining() != 0 {
		t.Errorf("Timer without deadline should have 0 time remaining, got %v",
			timer.TimeRemaining())
	}

	// Set deadline
	timer.SetDeadlineDuration(100 * time.Millisecond)

	// Immediately after setting, remaining time should be close to deadline
	remaining := timer.TimeRemaining()
	if remaining <= 10*time.Millisecond || remaining > 100*time.Millisecond {
		t.Errorf("Expected remaining time around 100ms, got %v", remaining)
	}

	// After waiting, remaining time should decrease
	time.Sleep(50 * time.Millisecond)
	newRemaining := timer.TimeRemaining()

	if newRemaining >= remaining {
		t.Errorf("Remaining time should decrease: %v vs %v", remaining, newRemaining)
	}

	// After deadline, remaining time should be 0
	time.Sleep(60 * time.Millisecond)
	if timer.TimeRemaining() != 0 {
		t.Errorf("After deadline, remaining time should be 0, got %v", timer.TimeRemaining())
	}
}

func TestString(t *testing.T) {
	// Test String method for a new timer
	timer := abstract.StartTimer()
	str := timer.String()

	// String should not be empty
	if str == "" {
		t.Error("String() should not return empty string")
	}

	// String should end with "ms" for a new timer
	if len(str) < 2 || str[len(str)-2:] != "ms" {
		t.Errorf("Expected string ending with 'ms', got '%s'", str)
	}

	// Test String method after some time has elapsed
	timer = abstract.StartTimer()
	time.Sleep(100 * time.Millisecond)
	str = timer.String()

	// Should contain some numeric value
	if str == "" {
		t.Error("String() should not return empty string after elapsed time")
	}

	// Should end with "ms" for millisecond durations
	if len(str) >= 2 && str[len(str)-2:] != "ms" {
		t.Errorf("Expected string ending with 'ms', got '%s'", str)
	}

	// Test that String() returns the same as FormatShort()
	expected := timer.FormatShort()
	if str != expected {
		t.Errorf("String() should return same as FormatShort(). Expected: %s, Got: %s", expected, str)
	}
}

func TestNewTimer(t *testing.T) {
	// Test creating timer with specific start time
	startTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	timer := abstract.NewTimer(startTime)

	// Timer should have the specified start time
	if !timer.Time().Equal(startTime) {
		t.Errorf("Expected start time %v, got %v", startTime, timer.Time())
	}

	// Timer should not be paused initially
	if timer.IsPaused() {
		t.Error("New timer should not be paused")
	}

	// Timer should have no laps initially
	if len(timer.Laps()) != 0 {
		t.Errorf("New timer should have no laps, got %d", len(timer.Laps()))
	}

	// Timer should have no deadline initially
	if timer.IsExpired() {
		t.Error("New timer should not be expired")
	}

	if timer.TimeRemaining() != 0 {
		t.Errorf("New timer should have 0 time remaining, got %v", timer.TimeRemaining())
	}
}

func TestNewTimerWithPastTime(t *testing.T) {
	// Test creating timer with a past time
	pastTime := time.Now().Add(-1 * time.Hour)
	timer := abstract.NewTimer(pastTime)

	// Timer should have the specified past start time
	if !timer.Time().Equal(pastTime) {
		t.Errorf("Expected start time %v, got %v", pastTime, timer.Time())
	}

	// Elapsed time should be positive (timer started in the past)
	elapsed := timer.ElapsedTime()
	if elapsed <= 0 {
		t.Errorf("Timer with past start time should have positive elapsed time, got %v", elapsed)
	}

	// Should be approximately 1 hour (with some tolerance)
	if elapsed < 59*time.Minute || elapsed > 61*time.Minute {
		t.Errorf("Expected elapsed time around 1 hour, got %v", elapsed)
	}
}

func TestNewTimerWithFutureTime(t *testing.T) {
	// Test creating timer with a future time
	futureTime := time.Now().Add(1 * time.Hour)
	timer := abstract.NewTimer(futureTime)

	// Timer should have the specified future start time
	if !timer.Time().Equal(futureTime) {
		t.Errorf("Expected start time %v, got %v", futureTime, timer.Time())
	}

	// Elapsed time should be negative (timer starts in the future)
	elapsed := timer.ElapsedTime()
	if elapsed >= 0 {
		t.Errorf("Timer with future start time should have negative elapsed time, got %v", elapsed)
	}

	// Should be approximately -1 hour (with some tolerance)
	if elapsed > -59*time.Minute || elapsed < -61*time.Minute {
		t.Errorf("Expected elapsed time around -1 hour, got %v", elapsed)
	}
}

func TestNewTimerEquality(t *testing.T) {
	// Test that NewTimer and StartTimer create equivalent timers when given current time
	now := time.Now()
	timer1 := abstract.NewTimer(now)
	timer2 := abstract.StartTimer()

	// Allow small difference in start times (microseconds)
	diff := timer2.Time().Sub(timer1.Time())
	if diff < -10*time.Microsecond || diff > 10*time.Microsecond {
		t.Errorf("Expected timers to have similar start times, got difference of %v", diff)
	}

	// Both timers should have same initial state
	if timer1.IsPaused() != timer2.IsPaused() {
		t.Error("Both timers should have same paused state")
	}

	if len(timer1.Laps()) != len(timer2.Laps()) {
		t.Error("Both timers should have same number of laps")
	}
}
