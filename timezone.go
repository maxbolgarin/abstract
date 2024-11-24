package abstract

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Timezone is a data structure to store timezone in UTC(+|-)HH:MM format.
type Timezone struct {
	loc    *time.Location
	offset int
}

// NewTimezone returns Timezone from provided [time.Location].
func NewTimezone(loc *time.Location) Timezone {
	if loc == nil {
		loc = time.UTC
	}
	return NewTimezoneFromTime(time.Now().In(loc))
}

// NewTimezoneFromTime returns Timezone from provided [time.Time].
func NewTimezoneFromTime(t time.Time) Timezone {
	_, offset := t.Zone()
	out := Timezone{
		offset: offset,
	}

	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}

	hours := offset / 3600
	minutes := offset % 3600 / 60
	if hours == 0 {
		out.loc = time.FixedZone("UTC", offset)
		return out
	}
	if minutes == 0 {
		out.loc = time.FixedZone(fmt.Sprintf("UTC%s%02d", sign, hours), offset)
		return out
	}

	out.loc = time.FixedZone(fmt.Sprintf("UTC%s%02d:%02d", sign, hours, minutes), offset)

	return out
}

// ParseTimezone returns Timezone from provided string - location or UTC(+|-)HH:MM.
func ParseTimezone(s string) (Timezone, error) {
	if len(s) < 3 {
		return Timezone{}, fmt.Errorf("invalid timezone: %s", s)
	}

	loc, err := time.LoadLocation(s)
	if err != nil {
		loc, err = ParseUTCOffset(s)
		if err != nil {
			return Timezone{}, err
		}
	}

	return NewTimezone(loc), nil
}

// Loc returns [time.Location] associated with Timezone.
func (i Timezone) Loc() *time.Location {
	return i.loc
}

// Offset returns offset in seconds.
func (i Timezone) Offset() int {
	return i.offset
}

// OffsetHours returns offset in hours.
func (i Timezone) OffsetHours() int {
	return i.offset / 3600
}

// String returns string representation of Timezone in UTC(+|-)HH:MM format.
func (i Timezone) String() string {
	return i.loc.String()
}

// MarshalJSON implements json.Marshaler interface to marshal Timezone to JSON.
func (i Timezone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + i.String() + `"`), nil
}

// MarshalYAML implements yaml.Marshaler interface to marshal Timezone to YAML.
func (i Timezone) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface to unmarshal Timezone from JSON.
func (i *Timezone) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	loc, err := ParseTimezone(s)
	if err != nil {
		return err
	}
	i.loc = loc.loc

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface to unmarshal Timezone from YAML.
func (i *Timezone) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	loc, err := ParseTimezone(s)
	if err != nil {
		return err
	}
	i.loc = loc.loc

	return nil
}

// ParseUTCOffset returns [time.Location] from provided string in UTC(+|-)HH:MM format.
func ParseUTCOffset(input string) (*time.Location, error) {
	if len(input) == 0 {
		return nil, errors.New("input cannot be empty")
	}
	input = strings.TrimSpace(strings.Replace(input, "UTC", "", 1))

	sign := byte('+')
	if input[0] == '+' || input[0] == '-' {
		if len(input) == 1 {
			return nil, errors.New("invalid input: " + input)
		}
		sign = input[0]
		input = input[1:]
	}

	var hours, minutes string

	for _, sep := range []string{" ", ":"} {
		spl := strings.Split(input, sep)
		if len(spl) == 2 {
			hours = spl[0]
			minutes = spl[1]
			break
		} else if len(spl) > 2 {
			return nil, errors.New("invalid input: " + input)
		}
		hours = spl[0]
		minutes = "0"
	}

	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		return nil, fmt.Errorf("invalid input %s and hours %s", input, hours)
	}
	hoursThreshold := 14
	if sign == '-' {
		hoursThreshold = 12
	}
	if hoursInt > hoursThreshold {
		return nil, fmt.Errorf("hours should be less than %d: %s", hoursThreshold, hours)
	}
	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		return nil, fmt.Errorf("invalid input %s and minutes %s", input, minutes)
	}
	if !isEqual(minutesInt, 0, 30, 45) {
		return nil, fmt.Errorf("minutes can be equal to 0, 30 or 45, got: %d", minutesInt)
	}
	if minutesInt == 30 {
		if sign == '+' {
			if !isEqual(hoursInt, 3, 4, 5, 6, 9, 10) {
				return nil, fmt.Errorf("invalid hour %s%d for minute %d", string(sign), hoursInt, minutesInt)
			}
		}
		if sign == '-' {
			if !isEqual(hoursInt, 3, 9) {
				return nil, fmt.Errorf("invalid hour %s%d for minute %d", string(sign), hoursInt, minutesInt)
			}
		}
	}
	if minutesInt == 45 {
		if sign == '-' {
			return nil, fmt.Errorf("invalid hour %s%d for minute %d", string(sign), hoursInt, minutesInt)
		}
		if sign == '+' {
			if !isEqual(hoursInt, 5, 8, 12) {
				return nil, fmt.Errorf("invalid hour %s%d for minute %d", string(sign), hoursInt, minutesInt)
			}
		}

	}
	signInt := 1
	if sign == '-' {
		signInt = -1
	}

	loc := strings.Builder{}
	loc.WriteString("UTC")
	loc.WriteByte(sign)
	loc.WriteString(hours)
	if minutesInt > 0 {
		loc.WriteString(":" + minutes)
	}
	return time.FixedZone(loc.String(), signInt*hoursInt*60*60+signInt*minutesInt*60), nil
}

func isEqual(n int, ns ...int) bool {
	for _, target := range ns {
		if n == target {
			return true
		}
	}
	return false
}
