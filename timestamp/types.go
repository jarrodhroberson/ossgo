package timestamp

import (
	"regexp"
	"strconv"
	"time"
)

const compact_format = "20060102t150405Z"

var iso8601DurationRegex = regexp.MustCompile(`^P(?:(?P<years>(\d+))Y)?(?:(?P<months>(\d+))M)?(?:(?P<weeks>(\d+))W)?(?:(?P<days>(\d+))D)?(?:T(?:(?P<hours>(\d+))H)?(?:(?P<minutes>(\d+))M)?(?:(?P<seconds>(\d+(?:\.\d+)?))S)?)?$`)

// Timestamps interface defines methods to access special timestamp values
type Timestamps interface {
	BeginningOfTime() *Timestamp
	EndOfTime() *Timestamp
	ZeroValue() *Timestamp
}

// enums implements the Timestamps interface and stores special timestamp values
type enums struct {
	beginningOfTime *Timestamp
	endOfTime       *Timestamp
	zeroValue       *Timestamp
}

// BeginningOfTime returns the earliest possible timestamp value
func (i enums) BeginningOfTime() *Timestamp {
	return i.beginningOfTime
}

// EndOfTime returns the latest possible timestamp value
func (i enums) EndOfTime() *Timestamp {
	return i.endOfTime
}

// ZeroValue returns the zero/uninitialized timestamp value
func (i enums) ZeroValue() *Timestamp {
	return i.zeroValue
}

// Timestamp wraps time.Time to provide additional functionality
type Timestamp struct {
	t time.Time
}

// Year returns the year of the timestamp
func (ts *Timestamp) Year() int {
	return ts.t.Year()
}

// Month returns the month of the timestamp
func (ts *Timestamp) Month() time.Month {
	return ts.t.Month()
}

// Day returns the weekday of the timestamp
func (ts *Timestamp) Day() time.Weekday {
	return ts.t.Weekday()
}

// Before reports whether the timestamp instant is before ots
func (ts *Timestamp) Before(ots *Timestamp) bool {
	return ts.t.Before(ots.t)
}

// After reports whether the timestamp instant is after ots
func (ts *Timestamp) After(ots *Timestamp) bool {
	return ts.t.After(ots.t)
}

// Add returns the timestamp t+d
func (ts *Timestamp) Add(d time.Duration) *Timestamp {
	return From(ts.t.Add(d))
}

// Sub returns the duration t-ots
func (ts *Timestamp) Sub(ots *Timestamp) time.Duration {
	return ts.t.Sub(ots.t)
}

// ZeroTime returns a new timestamp with the time portion set to midnight UTC
func (ts *Timestamp) ZeroTime() *Timestamp {
	return From(ts.t.Truncate(time.Hour * 24))
}

// Compare compares timestamps, returns -1 if t < to, 0 if t == to, 1 if t > to
func (ts *Timestamp) Compare(to *Timestamp) int {
	return ts.t.Compare(to.t)
}

// String returns the timestamp in RFC3339Nano format
func (ts *Timestamp) String() string {
	bytes, _ := ts.MarshalText()
	return string(bytes)
}

func (ts *Timestamp) MarshalText() (text []byte, err error) {
	return []byte(ts.t.UTC().Format(time.RFC3339Nano)), nil
}

func (ts *Timestamp) UnmarshalText(b []byte) error {
	t, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		return err
	}
	ts.t = t.UTC()
	return nil
}

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	bytes, err := ts.MarshalText()
	return []byte(strconv.Quote(string(bytes))), err
}

func (ts *Timestamp) UnmarshallJSON(b []byte) error {
	t, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		return err
	}
	ts.t = t.UTC()
	return nil
}

func (ts *Timestamp) MarshalBinary() (data []byte, err error) {
	return ts.MarshalText()
}

func (ts *Timestamp) UnmarshalBinary(b []byte) error {
	return ts.UnmarshalText(b)
}
