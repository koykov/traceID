package traceID

import (
	. "github.com/koykov/entry"
)

type EntryType uint8

const (
	EntryLog EntryType = iota
	EntryChapter
	EntryAcquireThread
	EntryReleaseThread
)

type entry struct {
	ll   LogLevel
	tp   EntryType
	tt   int64
	tid  uint32
	k, v Entry64
}
