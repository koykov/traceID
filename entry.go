package traceID

import (
	. "github.com/koykov/entry"
)

type EntryType uint8

const (
	EntryLog EntryType = iota
	EntrySubject
	EntryAcquireThread
	EntryReleaseThread
)

type entry struct {
	tp   EntryType
	tt   int64
	tid  uint32
	k, v Entry64
}
