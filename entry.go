package traceID

import (
	. "github.com/koykov/entry"
)

type EntryType uint8

const (
	EntryLog EntryType = iota
	EntrySubject
)

type entry struct {
	tp   EntryType
	tt   int64
	k, v Entry64
}
