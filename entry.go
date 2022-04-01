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
	ll   Level
	tp   EntryType
	tt   int64
	tid  uint32
	rid  uint32
	k, v Entry64
}

// Deferred entry.
type dentry struct {
	ll  Level
	tp  EntryType
	tt  int64
	tid uint32
	rid uint32
	k   string
	v   interface{}
	opt []optionKV
}

func (e EntryType) String() string {
	switch e {
	case EntryLog:
		return "LOG"
	case EntryChapter:
		return "CHAPTER"
	case EntryAcquireThread:
		return "THREAD_ACQ"
	case EntryReleaseThread:
		return "THREAD_REL"
	default:
		return "UNK"
	}
}
