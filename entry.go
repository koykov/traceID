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

func (d *dentry) reset() {
	d.ll = 0
	d.tp = 0
	d.tt = 0
	d.tid = 0
	d.rid = 0
	d.k = ""
	d.v = nil
	d.opt = d.opt[:0]
}

func (e EntryType) String() string {
	switch e {
	case EntryLog:
		return "LOG"
	case EntryChapter:
		return "CHAP"
	case EntryAcquireThread:
		return "TH_ACQ"
	case EntryReleaseThread:
		return "TH_REL"
	default:
		return "UNK"
	}
}
