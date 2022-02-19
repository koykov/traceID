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

func (c *Ctx) Debug(msg string) CtxInterface {
	c.chapter(LevelDebug, msg)
	return c
}

func (c *Ctx) Info(msg string) CtxInterface {
	c.chapter(LevelInfo, msg)
	return c
}

func (c *Ctx) Warn(msg string) CtxInterface {
	c.chapter(LevelWarn, msg)
	return c
}

func (c *Ctx) Error(msg string) CtxInterface {
	c.chapter(LevelError, msg)
	return c
}

func (c *Ctx) Fatal(msg string) CtxInterface {
	c.chapter(LevelFatal, msg)
	return c
}

func (c *Ctx) Var(name string, val interface{}) CtxInterface {
	c.log(LevelDebug, name, val, nil, false, EntryLog, 0)
	return c
}

func (c *Ctx) VarWithOptions(name string, val interface{}, opts Options) CtxInterface {
	c.log(LevelDebug, name, val, opts.Marshaller, opts.Indent, EntryLog, 0)
	return c
}

func (c *Ctx) VarIf(cond bool, name string, val interface{}) CtxInterface {
	if !cond {
		return c
	}
	c.log(LevelDebug, name, val, nil, false, EntryLog, 0)
	return c
}

func (c *Ctx) VarWithOptionsIf(cond bool, name string, val interface{}, opts Options) CtxInterface {
	if !cond {
		return c
	}
	c.log(LevelDebug, name, val, opts.Marshaller, opts.Indent, EntryLog, 0)
	return c
}

func (c *Ctx) chapter(level LogLevel, msg string) {
	c.log(level, "", msg, nil, false, EntryChapter, 0)
}

func (c *Ctx) log(level LogLevel, name string, val interface{}, m Marshaller, ind bool, typ EntryType, tid uint32) {
	c.lock()
	c.logLF(level, name, val, m, ind, typ, tid)
	c.unlock()
}

func (c *Ctx) logLF(level LogLevel, name string, val interface{}, m Marshaller, ind bool, typ EntryType, tid uint32) {
	off := len(c.buf)
	var k Entry64
	if l := len(name); l > 0 {
		c.buf = append(c.buf, name...)
		k.Encode(uint32(off), uint32(off+l))
	}

	off = len(c.buf)
	var v Entry64
	c.bb.Reset()
	if vb, err := c.getm(m).Marshal(&c.bb, val, ind); err == nil {
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
		ll:  level,
		tp:  typ,
		tt:  tt,
		tid: tid,
		k:   k,
		v:   v,
	})
}

func (c *Ctx) Flush() (err error) {
	c.lock()
	message := Encode(c)
	c.resetBuf()
	c.unlock()
	err = BroadcastWithTimeout(message, c.bto)
	return
}

func (c *Ctx) AcquireThread() ThreadInterface {
	id := atomic.AddUint32(&c.thc, 1)
	c.log(LevelDebug, "", id, nil, false, EntryAcquireThread, 0)
	return &Thread{
		id: id,
		rt: 0,
		cp: uintptr(unsafe.Pointer(c)),
	}
}

func (c *Ctx) ReleaseThread(thread ThreadInterface) CtxInterface {
	c.log(LevelDebug, "", thread.GetID(), nil, false, EntryReleaseThread, 0)
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
	sz += 2                                   // Version
	sz += 2                                   // ID length
	sz += len(c.id)                           // ID body
	sz += 2                                   // Entries count
	sz += len(c.lb) * (1 + 1 + 8 + 4 + 8 + 8) // Entry log level + type + timestamp + threadID + name + value
	sz += 4                                   // Payload length
	sz += len(c.buf)                          // Payload body
	return
}

func (c *Ctx) lock() {
	c.mux.Lock()
}

func (c *Ctx) unlock() {
	c.mux.Unlock()
}
