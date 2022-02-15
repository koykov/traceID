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

	bto time.Duration
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

func (c *Ctx) SetBroadcastTimeout(timeout time.Duration) CtxInterface {
	c.bto = timeout
	return c
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

func (c *Ctx) Subject(subject string) CtxInterface {
	c.log("", subject, nil, EntrySubject, 0)
	return c
}

func (c *Ctx) Log(key string, val interface{}) CtxInterface {
	c.log(key, val, nil, EntryLog, 0)
	return c
}

func (c *Ctx) LogWM(key string, val interface{}, m Marshaller) CtxInterface {
	c.log(key, val, m, EntryLog, 0)
	return c
}

func (c *Ctx) log(key string, val interface{}, m Marshaller, typ EntryType, tid uint32) {
	c.lock()
	defer c.unlock()

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

	var tt int64
	if c.cl != nil {
		tt = c.cl.Now().UnixNano()
	} else {
		tt = time.Now().UnixNano()
	}
	c.lb = append(c.lb, entry{
		tp:  typ,
		tt:  tt,
		tid: tid,
		k:   k,
		v:   v,
	})
}

func (c *Ctx) Commit() (err error) {
	c.lock()
	message := Encode(c)
	c.resetBuf()
	c.unlock()
	err = BroadcastWithTimeout(message, c.bto)
	return
}

func (c *Ctx) AcquireThread() ThreadInterface {
	id := atomic.AddUint32(&c.thc, 1)
	c.log("", id, nil, EntryAcquireThread, 0)
	return &Thread{
		id: id,
		rt: 0,
		cp: uintptr(unsafe.Pointer(c)),
	}
}

func (c *Ctx) ReleaseThread(thread ThreadInterface) CtxInterface {
	c.log("", thread.GetID(), nil, EntryReleaseThread, 0)
	return c
}

func (c *Ctx) Reset() *Ctx {
	c.bto = 0
	c.thc = 0
	c.resetBuf()
	c.m = nil
	c.bb.Reset()
	return c
}

func (c *Ctx) resetBuf() {
	c.lb = c.lb[:0]
	c.buf = c.buf[:0]
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
	sz += 2                               // Version
	sz += 2                               // ID length
	sz += len(c.id)                       // ID body
	sz += 2                               // Entries count
	sz += len(c.lb) * (1 + 8 + 4 + 8 + 8) // Entry type + timestamp + threadID + key + value
	sz += 4                               // Payload length
	sz += len(c.buf)                      // Payload body
	return
}

func (c *Ctx) lock() {
	c.mux.Lock()
}

func (c *Ctx) unlock() {
	c.mux.Unlock()
}
