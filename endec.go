package traceID

import (
	"encoding/binary"
	"fmt"

	"github.com/koykov/bytealg"
	. "github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

type MessageRow struct {
	Type  EntryType
	Time  int64
	Key   Entry64
	Value Entry64
}

type Message struct {
	Version uint16
	ID      string
	Rows    []MessageRow
	Buf     []byte
}

func Encode(ctx *Ctx) []byte {
	ctx.lock()

	poff := len(ctx.buf)
	size := ctx.size()
	ctx.buf = bytealg.GrowDelta(ctx.buf, size)
	buf := ctx.buf[poff:]
	off := 0
	binary.LittleEndian.PutUint16(buf[off:], Version)
	off += 2
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.id)))
	off += 2
	copy(buf[off:], ctx.id)
	off += len(ctx.id)
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.lb)))
	off += 2
	for i := 0; i < len(ctx.lb); i++ {
		e := &ctx.lb[i]
		buf[off] = uint8(e.tp)
		off++
		binary.LittleEndian.PutUint64(buf[off:], uint64(e.tt))
		off += 8
		binary.LittleEndian.PutUint64(buf[off:], uint64(e.k))
		off += 8
		binary.LittleEndian.PutUint64(buf[off:], uint64(e.v))
		off += 8
	}
	binary.LittleEndian.PutUint32(buf[off:], uint32(poff))
	off += 4
	copy(buf[off:], ctx.buf[:poff])

	ctx.unlock()
	return ctx.buf[poff:]
}

func Decode(p []byte, x *Message) error {
	if len(p) < 32 {
		return ErrPacketTooShort
	}
	off := 0
	if x.Version = binary.LittleEndian.Uint16(p[off:]); x.Version != Version {
		return fmt.Errorf("version mismatch: need %d, got %d", Version, x.Version)
	}
	off += 2
	l := binary.LittleEndian.Uint16(p[off:])
	off += 2
	if l >= uint16(len(p[off:])) {
		return ErrPacketTooShort
	}
	x.ID = fastconv.B2S(p[off : off+int(l)])
	off += int(l)
	if len(p[off:]) < 2 {
		return ErrPacketTooShort
	}
	if l = binary.LittleEndian.Uint16(p[off:]); l > 0 {
		x.Rows = make([]MessageRow, 0, l)
	}
	off += 2
	for i := 0; i < int(l); i++ {
		if len(p[off:]) < 25 {
			return ErrPacketTooShort
		}
		tp := EntryType(p[off])
		off++
		tt := binary.LittleEndian.Uint64(p[off:])
		off += 8
		k := binary.LittleEndian.Uint64(p[off:])
		off += 8
		v := binary.LittleEndian.Uint64(p[off:])
		off += 8
		x.Rows = append(x.Rows, MessageRow{
			Type:  tp,
			Time:  int64(tt),
			Key:   Entry64(k),
			Value: Entry64(v),
		})
	}
	if len(p[off:]) < 4 {
		return ErrPacketTooShort
	}
	bl := binary.LittleEndian.Uint32(p[off:])
	off += 4
	if bl > uint32(len(p[off:])) {
		return ErrPacketTooShort
	}
	x.Buf = p[off:]
	return nil
}
