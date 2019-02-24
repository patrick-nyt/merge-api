package main

import (
	"fmt"
	"time"
)

// JSONTime is used so we can have a custom time format
// when marshaling the timestamp into json.
type JSONTime struct {
	time.Time
}

// timeFormat is a custom time format. Not sure whether this
// was needed--- but since it was simple just put it in.
const timeFormat = "Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"

// MarshalJSON formats the timestamp using the custom format.
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format(timeFormat))
	return []byte(stamp), nil
}
