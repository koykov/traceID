package traceID

import (
	"time"

	"github.com/koykov/byteptr"
)

type EntryType uint8

const (
	EntryLog EntryType = iota
	EntrySubject
)

type entry struct {
	tp   EntryType
	tt   time.Time
	k, v byteptr.Byteptr
}
