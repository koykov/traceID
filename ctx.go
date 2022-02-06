package traceID

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/koykov/entry"
)

type Ctx struct {
	id string

	mux sync.Mutex
	lck uint32
	log []entry
	lb  []byte
	m   Marshaller
	c   Clock
	bb  bytes.Buffer
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
}

func (c *Ctx) SetClock(cl Clock) Interface {
	c.c = cl
	return c
}

func (c *Ctx) SetMarshaller(m Marshaller) Interface {
	c.m = m
	return c
}

func (c *Ctx) SetID(id string) Interface {
	c.id = id
	return c
}

func (c *Ctx) Subject(subject string) Interface {
	c._log("", subject, nil, EntrySubject)
	return c
}

func (c *Ctx) Log(key string, val interface{}) Interface {
	c._log(key, val, nil, EntryLog)
	return c
}

func (c *Ctx) LogWM(key string, val interface{}, m Marshaller) Interface {
	c._log(key, val, m, EntryLog)
	return c
}

func (c *Ctx) _log(key string, val interface{}, m Marshaller, typ EntryType) {
	off := len(c.lb)
	var k Entry64
	if l := len(key); l > 0 {
		c.lb = append(c.lb, key...)
		k.Encode(uint32(off), uint32(off+l))
	}

	off = len(c.lb)
	var v Entry64
	c.bb.Reset()
	if vb, err := c.getm(m).Marshal(&c.bb, val); err == nil {
		c.lb = append(c.lb, vb...)
		v.Encode(uint32(off), uint32(off+len(vb)))
	}

	var tt time.Time
	if c.c != nil {
		tt = c.c.Now()
	} else {
		tt = time.Now()
	}
	c.log = append(c.log, entry{
		tp: typ,
		tt: tt.UnixNano(),
		k:  k,
		v:  v,
	})
}

func (c *Ctx) BeginTXN() Interface {
	c.lock()
	c.mux.Lock()
	return c
}

func (c *Ctx) Commit() error {
	if !c.locked() {
		c.mux.Unlock()
		c.unlock()
	}
	return nil
}

func (c *Ctx) Reset() *Ctx {
	c.log = c.log[:0]
	c.lb = c.lb[:0]
	c.m = nil
	c.bb.Reset()
	if c.locked() {
		c.mux.Unlock()
		c.unlock()
	}
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
	sz += 2                            // Version
	sz += 2                            // ID length
	sz += len(c.id)                    // ID body
	sz += 2                            // Entries count
	sz += len(c.log) * (1 + 8 + 8 + 8) // Type + timestamp + key + value
	sz += 4                            // Payload length
	sz += len(c.lb)                    // Payload body
	return
}

func (c *Ctx) lock() {
	atomic.StoreUint32(&c.lck, 1)
}

func (c *Ctx) unlock() {
	atomic.StoreUint32(&c.lck, 0)
}

func (c *Ctx) locked() bool {
	return atomic.LoadUint32(&c.lck) == 1
}
