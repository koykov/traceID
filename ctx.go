package traceID

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/koykov/bitset"
	. "github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

type Ctx struct {
	id, svc string
	lmask   LogLevel

	bit bitset.Bitset
	bto time.Duration
	thc uint32
	rc  uint32
	mux sync.Mutex
	lb  []entry
	buf []byte
	m   Marshaller
	l   Logger
	cl  Clock
	bb  bytes.Buffer
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
}

func (c *Ctx) SetID(id string) CtxInterface {
	c.id = id
	return c
}

func (c *Ctx) SetService(svc string) CtxInterface {
	c.svc = svc
	return c
}

func (c *Ctx) SetFlag(flag Flag, value bool) CtxInterface {
	c.bit.SetBit(int(flag), value)
	return c
}

func (c *Ctx) Watch(mask LogLevel) CtxInterface {
	c.lmask = mask
	return c
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

func (c *Ctx) SetLogger(l Logger) CtxInterface {
	c.l = l
	return c
}

func (c *Ctx) Debug(msg string) RecordInterface {
	r := c.record(LevelDebug, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (c *Ctx) Info(msg string) RecordInterface {
	r := c.record(LevelInfo, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (c *Ctx) Warn(msg string) RecordInterface {
	r := c.record(LevelWarn, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (c *Ctx) Error(msg string) RecordInterface {
	r := c.record(LevelError, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (c *Ctx) Fatal(msg string) RecordInterface {
	r := c.record(LevelFatal, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (c *Ctx) record(level LogLevel, msg string) *Record {
	r := c.newRecord(0)
	c.log(level, "", msg, nil, false, EntryChapter, 0, r.id)
	return r
}

func (c *Ctx) log(level LogLevel, name string, val interface{}, m Marshaller, ind bool, typ EntryType, tid, rid uint32) {
	c.lock()
	c.logLF(level, name, val, m, ind, typ, tid, rid)
	c.unlock()
}

func (c *Ctx) logLF(level LogLevel, name string, val interface{}, m Marshaller, ind bool, typ EntryType, tid, rid uint32) {
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
		rid: rid,
		k:   k,
		v:   v,
	})
	if c.l != nil {
		c.l.Printf("[%s/%s] thread %d; record %d; key '%s'; %s\n",
			level.String(), typ.String(), tid, rid, name, fastconv.B2S(c.buf[off:]))
	}
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
	t := c.newThread(0)
	c.log(LevelDebug, "", t.id, nil, false, EntryAcquireThread, 0, 0)
	return t
}

func (c *Ctx) ReleaseThread(thread ThreadInterface) CtxInterface {
	c.log(LevelDebug, "", thread.GetID(), nil, false, EntryReleaseThread, 0, 0)
	return c
}

func (c *Ctx) Reset() *Ctx {
	c.lmask = LogAll
	c.bit = 0
	c.bto = 0
	c.thc = 0
	c.rc = 0
	c.resetBuf()
	c.m = nil
	c.l = nil
	c.bb.Reset()
	return c
}

func (c *Ctx) resetBuf() {
	c.lb = c.lb[:0]
	c.buf = c.buf[:0]
}

func (c *Ctx) newThread(root uint32) *Thread {
	id := atomic.AddUint32(&c.thc, 1)
	t := &Thread{
		ctxHeir: ctxHeir{cp: uintptr(unsafe.Pointer(c))},
		id:      id,
		rt:      root,
	}
	return t
}

func (c *Ctx) newRecord(thid uint32) *Record {
	id := atomic.AddUint32(&c.rc, 1)
	r := &Record{
		ctxHeir: ctxHeir{cp: uintptr(unsafe.Pointer(c))},
		id:      id - 1,
		thid:    thid,
	}
	return r
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
	sz += 2                                       // Version
	sz += 8                                       // Bitset
	sz += 2                                       // ID length
	sz += len(c.id)                               // ID body
	sz += 2                                       // Service length
	sz += len(c.svc)                              // Service body
	sz += 2                                       // Entries count
	sz += len(c.lb) * (1 + 1 + 8 + 4 + 4 + 8 + 8) // Entry log level + type + timestamp + threadID + recordID + name + value
	sz += 4                                       // Payload length
	sz += len(c.buf)                              // Payload body
	return
}

func (c *Ctx) lock() {
	c.mux.Lock()
}

func (c *Ctx) unlock() {
	c.mux.Unlock()
}
