package timestamp

import (
	"fmt"
	"time"
)

// Period represents a time period with a start and end timestamp.
type Period struct {
	Start *Timestamp `json:"start"`
	End   *Timestamp `json:"end"`
}

// String returns a string representation of the Period in the format "start|end".
func (p *Period) String() string {
	return fmt.Sprintf("%s|%s", p.Start, p.End)
}

// Duration returns the duration of the Period.
// It calculates the absolute difference between the end and start timestamps.
func (p *Period) Duration() time.Duration {
	return p.End.Sub(p.Start).Abs()
}

// Contains checks if the given Timestamp is within the Period (inclusive of start and end).
func (p *Period) Contains(ts *Timestamp) bool {
	return (p.Start == ts || p.End == ts) || (p.Start.Before(ts) && p.End.After(ts))
}
