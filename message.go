package traceID

import (
	"github.com/koykov/bitset"
	. "github.com/koykov/entry"
)

type MessageRow struct {
	Level    LogLevel
	Type     EntryType
	Time     int64
	ThreadID uint32
	RecordID uint32
	Key      Entry64
	Value    Entry64
}

type Message struct {
	Bits    bitset.Bitset
	Version uint16
	ID      string
	Service string
	Rows    []MessageRow
	Buf     []byte
}

func (m Message) CheckFlag(flag Flag) bool {
	return m.Bits.CheckBit(int(flag))
}
