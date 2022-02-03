package traceID

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/byteptr"
)

type Ctx struct {
	id string

	mux sync.Mutex
	lck uint32
	log []entry
	lb  []byte
	m   Marshaller
	bb  bytes.Buffer
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
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
	c._log("", subject, EntrySubject)
	return c
}

func (c *Ctx) Log(key string, val interface{}) Interface {
	c._log(key, val, EntryLog)
	return c
}

func (c *Ctx) _log(key string, val interface{}, typ EntryType) {
	off := len(c.lb)
	k := byteptr.Byteptr{}
	if l := len(key); l > 0 {
		c.lb = append(c.lb, key...)
		k.Init(c.lb, off, l)
	}

	var m Marshaller
	if m = c.m; m == nil {
		m = defaultMarshaller
	}
	off = len(c.lb)
	v := byteptr.Byteptr{}
	c.bb.Reset()
	if vb, err := m.Marshal(&c.bb, val); err == nil {
		c.lb = append(c.lb, vb...)
		v.Init(c.lb, off, len(vb))
	}

	c.log = append(c.log, entry{
		tp: typ,
		tt: time.Now(),
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
	c.m = nil
	c.bb.Reset()
	return c
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
