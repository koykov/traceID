package traceID

import (
	"encoding/binary"

	"github.com/koykov/bytealg"
)

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
