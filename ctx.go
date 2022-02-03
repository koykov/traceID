package traceID

import (
	"sync"
	"time"
)

type Ctx struct {
	component, id string

	mux sync.Mutex
	log []entry
	ll  int
	m   Marshaller
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
	c.mux.Lock()
	return c
}

func (c *Ctx) Commit() error {
	c.mux.Unlock()
	return c.Push()
}

func (c *Ctx) Reset() *Ctx {
	for i := 0; i < c.ll; i++ {
		c.log[i].t = 0
		c.log[i].k = ""
		c.log[i].v = nil
	}
	c.ll = 0
	c.m = nil
	return c
}
