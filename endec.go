package traceID

import (
	"encoding/binary"
	"fmt"

	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	. "github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

func Encode(ctx *Ctx) []byte {
	ctx.flushDL()
	poff := len(ctx.buf)
	size := ctx.size()
	ctx.buf = bytealg.GrowDelta(ctx.buf, size)
	buf := ctx.buf[poff:]
	off := 0
	binary.LittleEndian.PutUint16(buf[off:], Version)
	off += 2
	binary.LittleEndian.PutUint64(buf[off:], uint64(ctx.bit))
	off += 8
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.id)))
	off += 2
	copy(buf[off:], ctx.id)
	off += len(ctx.id)
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.svc)))
	off += 2
	copy(buf[off:], ctx.svc)
	off += len(ctx.svc)
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.stg)))
	off += 2
	copy(buf[off:], ctx.stg)
	off += len(ctx.stg)
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.lb)))
	off += 2
	for i := 0; i < len(ctx.lb); i++ {
		e := &ctx.lb[i]
		buf[off] = uint8(e.ll)
		off++
		buf[off] = uint8(e.tp)
		off++
		binary.LittleEndian.PutUint64(buf[off:], uint64(e.tt))
		off += 8
		binary.LittleEndian.PutUint32(buf[off:], e.tid)
		off += 4
		binary.LittleEndian.PutUint32(buf[off:], e.rid)
		off += 4
		binary.LittleEndian.PutUint64(buf[off:], uint64(e.k))
		off += 8
		binary.LittleEndian.PutUint64(buf[off:], uint64(e.v))
		off += 8
	}
	binary.LittleEndian.PutUint32(buf[off:], uint32(poff))
	off += 4
	copy(buf[off:], ctx.buf[:poff])

	return ctx.buf[poff:]
}

func Decode(p []byte, x *Message) error {
	off := 0
	if len(p[off:]) < 2 {
		return ErrPacketTooShort
	}
	if x.Version = binary.LittleEndian.Uint16(p[off:]); x.Version != Version {
		return fmt.Errorf("version mismatch: need %d, got %d", Version, x.Version)
	}
	off += 2

	if len(p[off:]) < 8 {
		return ErrPacketTooShort
	}
	x.Bits = bitset.Bitset(binary.LittleEndian.Uint64(p[off:]))
	off += 8

	if len(p[off:]) < 2 {
		return ErrPacketTooShort
	}
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
	l = binary.LittleEndian.Uint16(p[off:])
	off += 2
	if l >= uint16(len(p[off:])) {
		return ErrPacketTooShort
	}
	x.Service = fastconv.B2S(p[off : off+int(l)])
	off += int(l)

	if len(p[off:]) < 2 {
		return ErrPacketTooShort
	}
	l = binary.LittleEndian.Uint16(p[off:])
	off += 2
	if l >= uint16(len(p[off:])) {
		return ErrPacketTooShort
	}
	x.Stage = fastconv.B2S(p[off : off+int(l)])
	off += int(l)

	if len(p[off:]) < 2 {
		return ErrPacketTooShort
	}
	if l = binary.LittleEndian.Uint16(p[off:]); l > 0 {
		x.Rows = make([]MessageRow, 0, l)
	}
	off += 2
	for i := 0; i < int(l); i++ {
		if len(p[off:]) < 30 {
			return ErrPacketTooShort
		}
		ll := Level(p[off])
		off++
		tp := EntryType(p[off])
		off++
		tt := binary.LittleEndian.Uint64(p[off:])
		off += 8
		tid := binary.LittleEndian.Uint32(p[off:])
		off += 4
		rid := binary.LittleEndian.Uint32(p[off:])
		off += 4
		k := binary.LittleEndian.Uint64(p[off:])
		off += 8
		v := binary.LittleEndian.Uint64(p[off:])
		off += 8
		x.Rows = append(x.Rows, MessageRow{
			Level:    ll,
			Type:     tp,
			Time:     int64(tt),
			ThreadID: tid,
			RecordID: rid,
			Key:      Entry64(k),
			Value:    Entry64(v),
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
