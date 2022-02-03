package traceID

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
)

type Ctx struct {
	component, id string

	mux sync.Mutex
	lck uint32
	log []entry
	ll  int
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

func (c *Ctx) SetComponent(component string) Interface {
	c.component = component
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
	if c.ll < len(c.log) {
		c.log[c.ll].t = typ
		c.log[c.ll].tt = time.Now()
		c.log[c.ll].k = key
		c.log[c.ll].v = val
	} else {
		c.log = append(c.log, entry{
			t:  typ,
			tt: time.Now(),
			k:  key,
			v:  val,
		})
	}
}

func (c *Ctx) Push() error {

	return nil
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
	if c.m == nil {
		return ErrNoMarshaller
	}
	return nil
}

func (c *Ctx) Reset() *Ctx {
	for i := 0; i < c.ll; i++ {
		c.log[i].t = 0
		c.log[i].k = ""
		c.log[i].v = nil
	}
	c.ll = 0
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
