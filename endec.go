package traceID

import (
	"encoding/binary"
	"fmt"

	"github.com/koykov/bytealg"
	. "github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

type Record struct {
	Type  EntryType
	Time  int64
	Key   Entry64
	Value Entry64
}

type Packet struct {
	Version uint16
	ID      string
	Records []Record
	Buf     []byte
}

func Encode(ctx *Ctx) []byte {
	ctx.lock()
	ctx.mux.Lock()

	poff := len(ctx.lb)
	size := ctx.size()
	ctx.lb = bytealg.GrowDelta(ctx.lb, size)
	buf := ctx.lb[poff:]
	off := 0
	binary.LittleEndian.PutUint16(buf[off:], Version)
	off += 2
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.id)))
	off += 2
	copy(buf[off:], ctx.id)
	off += len(ctx.id)
	binary.LittleEndian.PutUint16(buf[off:], uint16(len(ctx.log)))
	for i := 0; i < len(ctx.log); i++ {
		e := &ctx.log[i]
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
	copy(buf[off:], ctx.lb[:poff])

	ctx.mux.Unlock()
	ctx.unlock()
	return ctx.lb[poff:]
}

func Decode(p []byte, x *Packet) error {
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
		x.Records = make([]Record, 0, l)
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
		x.Records = append(x.Records, Record{
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
