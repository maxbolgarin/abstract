package abstract

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

// EmptyDate is a not initialized Date.
var EmptyDate = Date{}

// Date is a data structure to store date without time.
type Date struct {
	time.Time
}

// NewDate returns new date from year, month and day.
func NewDate(year, month, day int) Date {
	return Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
}

// NewDateFromString returns new date from yyyy-mm-dd string.
func NewDateFromString(date string) (Date, error) {
	d, err := time.Parse(dateLayout, date)
	if err != nil {
		return Date{}, err
	}
	return NewDateFromTime(d), nil
}

// NewDateFromTime returns new date from time.Time.
func NewDateFromTime(t time.Time) Date {
	return Date{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

// NowDate returns current active day.
func NowDate(tz *time.Location) Date {
	now := time.Now().In(tz)
	return NewDate(now.Year(), int(now.Month()), now.Day())
}

// Today returns current active day according to dayStart time.
func Today(dayStart Time, tz *time.Location) Date {
	now := time.Now().In(tz)
	if now.Hour() < dayStart.Hour() || (now.Hour() == dayStart.Hour() && now.Minute() < dayStart.Minute()) {
		now = now.AddDate(0, 0, -1)
	}
	return NewDate(now.Year(), int(now.Month()), now.Day())
}

// ParseDate tries to parse date (yyyy-mm-dd) using separators: ["-", " ", ".", "-", "_"].
func ParseDate(s string) (Date, error) {
	if s == "" {
		return Date{}, errors.New("date is empty")
	}
	seps := []string{"-", " ", ".", "-", "_"}
	for _, sep := range seps {
		splitted := strings.Split(s, sep)
		if len(splitted) == 3 {
			year, err := strconv.Atoi(splitted[0])
			if err != nil {
				return Date{}, fmt.Errorf("parse year=%s: %w", splitted[0], err)
			}

			month, err := strconv.Atoi(splitted[1])
			if err != nil {
				return Date{}, fmt.Errorf("parse month=%s: %w", splitted[1], err)
			}

			day, err := strconv.Atoi(splitted[2])
			if err != nil {
				return Date{}, fmt.Errorf("parse day=%s: %w", splitted[2], err)
			}

			return NewDate(year, month, day), nil
		}
	}
	return Date{}, fmt.Errorf("invalid date=%s", s)
}

// SortDates sorts dates.
func SortDates(dates []Date, desc bool) {
	sort.Slice(dates, func(i, j int) bool {
		if desc {
			return dates[i].After(dates[j].Time)
		}
		return dates[i].Before(dates[j].Time)
	})
}

// String returns date in yyyy-mm-dd format.
func (d Date) String() string {
	return d.Format(dateLayout)
}

// Round returns new Date instance with Round(0).
func (d Date) Round() Date {
	return Date{d.Time.Round(0)}
}

// NextDay returns Date instance for the next day.
func (d Date) NextDay() Date {
	d.Time = d.AddDate(0, 0, 1)
	return NewDateFromTime(d.Time)
}

// PrevDay returns Date instance for the previous day.
func (d Date) PrevDay() Date {
	d.Time = d.AddDate(0, 0, -1)
	return NewDateFromTime(d.Time)
}

// IsZero returns true if date is empty.
func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

// EqualDate returns true if dates are equal.
func (d Date) EqualDate(other Date) bool {
	return d.Day() == other.Day() && d.Month() == other.Month() && d.Year() == other.Year()
}

// Range returns number of days between two dates.
func (d Date) Range(other Date) int {
	d1 := d.Unix()
	d2 := other.Unix()
	r := int(d2 - d1)
	if r < 0 {
		r *= -1
	}
	return r / 86400
}

// IsToday returns true if provided argument is today.
func (d Date) IsToday(dayStart Time, tz *time.Location) bool {
	return d.EqualDate(Today(dayStart, tz))
}

// IsArgNextDay returns true if provided argument is after Date.
func (d Date) IsArgNextDay(t Date) bool {
	if d.Year() < t.Year() {
		return true
	} else if d.Year() > t.Year() {
		return false
	}

	if d.Month() < t.Month() {
		return true
	} else if d.Month() > t.Month() {
		return false
	}

	if d.Day() < t.Day() {
		return true
	} else if d.Day() > t.Day() {
		return false
	}

	return false
}

// MarshalJSON implements json.Marshaler interface to marshal Date to JSON.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// MarshalYAML implements yaml.Marshaler interface to marshal Date to YAML.
func (d Date) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface to unmarshal Date from JSON.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	res, err := NewDateFromString(s)
	if err != nil {
		return err
	}
	d.Time = res.Time
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface to unmarshal Date from YAML.
func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	res, err := ParseDate(s)
	if err != nil {
		return err
	}
	d.Time = res.Time

	return nil
}

// TransformDatesToString transforms slice of dates to slice of strings.
func TransformDatesToString(dates []Date) []string {
	out := make([]string, 0, len(dates))
	for _, d := range dates {
		out = append(out, d.String())
	}
	return out
}
