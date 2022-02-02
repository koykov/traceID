package traceID

import (
	"sync"
	"time"
)

type Ctx struct {
	c7t, id string

	mux sync.Mutex
	log []tuple
	ll  int
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
}

func (c *Ctx) SetID(component, id string) {
	c.c7t, c.id = component, id
}

func (c *Ctx) BeginTXN() Interface {
	c.Reset()
	return c
}

func (c *Ctx) Log(key string, val interface{}) Interface {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.ll < len(c.log) {
		c.log[c.ll].t = time.Now()
		c.log[c.ll].k = key
		c.log[c.ll].v = val
	} else {
		c.log = append(c.log, tuple{
			t: time.Now(),
			k: key,
			v: val,
		})
	}
	return c
}

func (c *Ctx) Commit() error {
	// ...
	return nil
}

func (c *Ctx) CommitTo(bc Broadcaster) error {
	// ...
	c.Reset()
	return nil
}

func (c *Ctx) Reset() *Ctx {
	for i := 0; i < c.ll; i++ {
		c.log[i].k = ""
		c.log[i].v = nil
	}
	return c
}
