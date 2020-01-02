package es

import "time"

// TODO: rename to Entry
type Record struct {
	At   time.Time
	Type string
	Data string
}
