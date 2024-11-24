package abstract

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	minutesInDay = 60 * 24
	secondsInDay = 86400
	timeLayout   = "15:04"
)

// EmptyTime is a not initialized Time.
var EmptyTime = Time{}

// Time is a data structure to store hours and minutes.
type Time struct {
	time.Time
	isSet bool
}

// NewTime returns new time from hour and minute.
func NewTime(hour, minute int) Time {
	return Time{time.Date(0, 0, 0, hour, minute, 0, 0, time.UTC), true}
}

// NewTimeFromString returns new time from HH:MM string.
func NewTimeFromString(s string) (Time, error) {
	d, err := time.Parse(timeLayout, s)
	if err != nil {
		return Time{}, err
	}
	return NewFromTime(d), nil
}

// NewFromTime returns new Time from time.Time.
func NewFromTime(t time.Time) Time {
	return Time{time.Date(0, 0, 0, t.Hour(), t.Minute(), 0, 0, time.UTC), true}
}

// NowTime returns current time.
func NowTime(tz *time.Location) Time {
	now := time.Now().In(tz)
	return NewFromTime(now)
}

// ParseTime tries to parse time (HH:MM) using separators: [" ", ":", "-", "_", ",", "."].
func ParseTime(s string) (Time, error) {
	if s == "" {
		return Time{}, errors.New("time is empty")
	}
	seps := []string{" ", ":", "-", "_", ",", "."}
	for _, sep := range seps {
		splitted := strings.Split(s, sep)
		if len(splitted) != 2 {
			if len(s) != 4 {
				continue
			}
			splitted = []string{string(s[0:2]), string(s[2:4])}
		}

		splitted[0] = prepareNumber(splitted[0], false)
		hour, err := strconv.Atoi(splitted[0])
		if err != nil {
			return Time{}, fmt.Errorf("parse hour=%s: %w", splitted[0], err)
		}
		if hour < 0 || hour > 23 {
			return Time{}, fmt.Errorf("invalid hour=%d", hour)
		}

		splitted[1] = prepareNumber(splitted[1], false)
		minute, err := strconv.Atoi(splitted[1])
		if err != nil {
			return Time{}, fmt.Errorf("parse minute=%s: %w", splitted[1], err)
		}
		if minute < 0 || minute > 59 {
			return Time{}, fmt.Errorf("invalid minute=%d", minute)
		}

		return NewTime(hour, minute), nil
	}

	return Time{}, fmt.Errorf("invalid time=%s", s)
}

// String returns time in HH:MM format.
func (t Time) String() string {
	return t.Format(timeLayout)
}

// Range substracts low from high time and returns duration between it.
func (low Time) Range(high Time) time.Duration {
	return time.Hour*time.Duration(high.Hour()-low.Hour()) +
		time.Minute*time.Duration(high.Minute()-low.Minute())
}

// RangeUp returns duration from low time to high time ignoring dates.
func (low Time) RangeUp(high Time) time.Duration {
	var hours, minutes int
	if high.Hour() < low.Hour() {
		hours = 24 - low.Hour() + high.Hour()
	} else {
		hours = high.Hour() - low.Hour()
	}

	if high.Minute() < low.Minute() {
		hours -= 1
		if hours < 0 {
			hours = 23
		}
		minutes = 60 - low.Minute() + high.Minute()
	} else {
		minutes = high.Minute() - low.Minute()
	}

	return time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes)
}

// AddTime adds howMuch to time.
func (t Time) AddTime(howMuch time.Duration) Time {
	minutes := int(howMuch.Minutes())
	for minutes > minutesInDay {
		minutes -= minutesInDay
	}
	tillEnd := t.MinutesTillDayEnd(EmptyTime)
	if minutes < tillEnd {
		return t.addTime(time.Minute, minutes)
	}
	afterDayStart := minutes - tillEnd
	return NewTime(0, 0).addTime(time.Minute, afterDayStart)
}

func (t Time) addTime(what time.Duration, howMuch int) Time {
	return Time{t.Add(what * time.Duration(howMuch)), true}
}

// SubTime substracts howMuch from time.
func (t Time) SubTime(howMuch time.Duration) Time {
	minutes := int(howMuch.Minutes())
	fromBegin := t.MinutesFromDayBegin(EmptyTime)
	for minutes > minutesInDay {
		minutes -= minutesInDay
	}
	if minutes < fromBegin {
		return t.subTime(time.Minute, minutes)
	}
	beforeDayStart := minutes - fromBegin
	return NewTime(23, 59).subTime(time.Minute, beforeDayStart+1)
}

func (t Time) subTime(what time.Duration, howMuch int) Time {
	return t.addTime(what, -1*howMuch)
}

// MinutesFromDayBegin returns number of minutes passed from the beginning of the day.
func (t Time) MinutesFromDayBegin(dayStartTime Time) int {
	var hours int
	if t.Hour() < dayStartTime.Hour() {
		hours = 24 - dayStartTime.Hour() - t.Hour()
	} else {
		hours = t.Hour() - dayStartTime.Hour()
	}
	return hours*60 + t.Minute() - 1
}

// MinutesTillDayEnd returns number of minutes remaining to the end of the day.
func (t Time) MinutesTillDayEnd(dayStartTime Time) int {
	return minutesInDay - t.MinutesFromDayBegin(dayStartTime)
}

// EqualTime returns true if times are equal.
func (t Time) EqualTime(other Time) bool {
	return t.Hour() == other.Hour() && t.Minute() == other.Minute()
}

// IsBefore returns true if reciever is before or equal to argument.
func (t Time) IsBefore(other Time) bool {
	if t.Hour() > other.Hour() {
		return false
	}
	if t.Hour() == other.Hour() && t.Minute() > other.Minute() {
		return false
	}
	return true
}

// IsArgBefore returns true if argument is before or equal to reciever.
func (t Time) IsArgBefore(other Time) bool {
	if t.Hour() < other.Hour() {
		return false
	}
	if t.Hour() == other.Hour() && t.Minute() < other.Minute() {
		return false
	}
	return true
}

// IsBefore returns true if reciever is STRICTLY before argument.
func (t Time) IsBeforeStrict(other Time) bool {
	if t.Hour() > other.Hour() {
		return false
	}
	if t.Hour() == other.Hour() && t.Minute() >= other.Minute() {
		return false
	}
	return true
}

// IsArgBefore returns true if argument is STRICTLY before reciever.
func (t Time) IsArgBeforeStrict(other Time) bool {
	if t.Hour() < other.Hour() {
		return false
	}
	if t.Hour() == other.Hour() && t.Minute() <= other.Minute() {
		return false
	}
	return true
}

// IsAfter returns true if reciever is after or equal to argument.
func (t Time) IsAfter(other Time) bool {
	return t.IsArgBefore(other)
}

// IsArgAfter returns true if argument is after or equal to reciever.
func (t Time) IsArgAfter(other Time) bool {
	return t.IsBefore(other)
}

// IsAfterStrict returns true if reciever is STRICTLY after argument.
func (t Time) IsAfterStrict(other Time) bool {
	return t.IsArgBeforeStrict(other)
}

// IsArgAfterStrict returns true if argument is STRICTLY after reciever.
func (t Time) IsArgAfterStrict(other Time) bool {
	return t.IsBeforeStrict(other)
}

// SmartDiff returns diff where reciever is start and argument is end
func (start Time) SmartDiff(end Time) time.Duration {
	var (
		startMinutes = start.MinutesFromDayBegin(EmptyTime)
		endMinutes   = end.MinutesFromDayBegin(EmptyTime)
	)

	if endMinutes >= startMinutes {
		return time.Minute * time.Duration(endMinutes-startMinutes)
	}
	return time.Minute * time.Duration(endMinutes+start.MinutesTillDayEnd(EmptyTime))
}

// RoundToFives returns time rounded to nearest 5 minutes
func (t Time) RoundToFives() Time {
	m := t.Minute()
	for i := 1; i <= 6; i++ {
		base := i * 10
		if m < base {
			if m <= base-5 {
				m = base - 10
			} else {
				m = base - 5
			}
			break
		}
	}
	return NewTime(t.Hour(), m)
}

// RoundToFivesUp adds 5 minutes and then RoundToFives
func (t Time) RoundToFivesUp() Time {
	return NewFromTime(t.RoundToFives().Add(5 * time.Minute))
}

// IsZero returns true if time is empty.
func (t Time) IsZero() bool {
	if t.Time.IsZero() {
		return !t.isSet
	}
	return false
}

// MarshalJSON implements json.Marshaler interface to marshal Time to JSON.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.isSet {
		return []byte("null"), nil
	}
	return json.Marshal(t.String())
}

// MarshalYAML implements yaml.Marshaler interface to marshal Time to YAML.
func (i Time) MarshalYAML() (interface{}, error) {
	if !i.isSet {
		return nil, nil
	}
	return i.String(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface to unmarshal Time from JSON.
func (i *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	res, err := ParseTime(s)
	if err != nil {
		return err
	}
	i.Time = res.Time
	i.isSet = true

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface to unmarshal Time from YAML.
func (e *Time) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var t string
	err := unmarshal(&t)
	if err != nil {
		return err
	}

	res, err := ParseTime(t)
	if err != nil {
		return err
	}
	*e = res

	return nil
}

func prepareNumber(s string, isDecimal bool) string {
	for i := range s {
		if s[i] >= '0' && s[i] <= '9' {
			continue
		}
		if isDecimal && s[i] == '.' {
			continue
		}
		return s[:i]
	}
	return s
}

// SortingPriority is a mark to group time for sorting.
type SortingPriority int

const (
	// LongAgoPriority shoud be first in the list because it should be handled a long time ago.
	LongAgoPriority SortingPriority = iota + 1
	// BeforePriority should be second in the list because it should be handled not long ago.
	BeforePriority
	// AfterPriority should be third in the list because it should be handled soon.
	AfterPriority
	// NotSoonPriority should be last in the list because it should be handled not soon.
	NotSoonPriority
)

// GetTimeSortingPriority returns time priority for sorting.
// toCheck is time to check, now is current time, dayStart is day start time (00:00 in general, but not always).
func GetTimeSortingPriority(toCheck, now, dayStart Time) SortingPriority {
	if dayStart.IsArgBefore(now) {
		if toCheck.IsBefore(dayStart) {
			if toCheck.IsBefore(now) {
				return BeforePriority
			}
			return AfterPriority
		}
		return LongAgoPriority
	}

	if toCheck.IsAfter(dayStart) {
		if toCheck.IsBefore(now) {
			return BeforePriority
		}
		return AfterPriority
	}

	return NotSoonPriority
}
