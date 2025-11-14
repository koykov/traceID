package traceID

import (
	"github.com/koykov/bitset"
	ent "github.com/koykov/entry"
)

type MessageRow struct {
	Level    Level
	Type     EntryType
	Time     int64
	ThreadID uint32
	RecordID uint32
	Key      ent.Entry64
	Value    ent.Entry64
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
