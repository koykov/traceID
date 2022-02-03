package traceID

import "time"

type EntryType uint8

const (
	EntryLog EntryType = iota
	EntrySubject
)

type entry struct {
	t  EntryType
	tt time.Time
	k  string
	v  interface{}
}
