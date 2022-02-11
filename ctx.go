package traceID

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	. "github.com/koykov/entry"
)

type Ctx struct {
	id string

	thc uint32
	mux sync.Mutex
	lb  []entry
	buf []byte
	m   Marshaller
	cl  Clock
	bb  bytes.Buffer
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
}

func (c *Ctx) SetClock(cl Clock) CtxInterface {
	c.cl = cl
	return c
}

func (c *Ctx) SetMarshaller(m Marshaller) CtxInterface {
	c.m = m
	return c
}

func (c *Ctx) SetID(id string) CtxInterface {
	c.id = id
	return c
}

func (c *Ctx) Thread() ThreadInterface {
	id := atomic.AddUint32(&c.thc, 1)
	return &Thread{
		id: id,
		cp: uintptr(unsafe.Pointer(c)),
	}
}

func (c *Ctx) Subject(subject string) CtxInterface {
	c.logLF("", subject, nil, EntrySubject)
	return c
}

func (c *Ctx) Log(key string, val interface{}) CtxInterface {
	c.logLF(key, val, nil, EntryLog)
	return c
}

func (c *Ctx) LogWM(key string, val interface{}, m Marshaller) CtxInterface {
	c.logLF(key, val, m, EntryLog)
	return c
}

func (c *Ctx) logLF(key string, val interface{}, m Marshaller, typ EntryType) {
	off := len(c.buf)
	var k Entry64
	if l := len(key); l > 0 {
		c.buf = append(c.buf, key...)
		k.Encode(uint32(off), uint32(off+l))
	}

	off = len(c.buf)
	var v Entry64
	c.bb.Reset()
	if vb, err := c.getm(m).Marshal(&c.bb, val); err == nil {
		c.buf = append(c.buf, vb...)
		v.Encode(uint32(off), uint32(off+len(vb)))
	}

	var tt time.Time
	if c.cl != nil {
		tt = c.cl.Now()
	} else {
		tt = time.Now()
	}
	c.lb = append(c.lb, entry{
		tp: typ,
		tt: tt.UnixNano(),
		k:  k,
		v:  v,
	})
}

func (c *Ctx) Commit() error {
	// ...
	return nil
}

func (c *Ctx) Reset() *Ctx {
	c.thc = 0
	c.lb = c.lb[:0]
	c.buf = c.buf[:0]
	c.m = nil
	c.bb.Reset()
	return c
}

func (c *Ctx) getm(m Marshaller) Marshaller {
	if m != nil {
		return m
	}
	if c.m != nil {
		return c.m
	}
	return defaultMarshaller
}

func (c *Ctx) size() (sz int) {
	sz += 2                           // Version
	sz += 2                           // ID length
	sz += len(c.id)                   // ID body
	sz += 2                           // Entries count
	sz += len(c.lb) * (1 + 8 + 8 + 8) // Type + timestamp + key + value
	sz += 4                           // Payload length
	sz += len(c.buf)                  // Payload body
	return
}

func (c *Ctx) lock() {
	c.mux.Lock()
}

func (c *Ctx) unlock() {
	c.mux.Unlock()
}
