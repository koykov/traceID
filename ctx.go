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
	"github.com/koykov/x2bytes"
)

type Ctx struct {
	id, svc string
	lmask   Level

	bit  bitset.Bitset
	bto  time.Duration
	thc  uint32
	rc   uint32
	mux  sync.Mutex
	lb   []entry
	dlb  []dentry
	dlbc int
	buf  []byte
	m    Marshaller
	l    Logger
	cl   Clock
	bb   bytes.Buffer
}

func NewCtx() *Ctx {
	ctx := Ctx{lmask: LevelAll}
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

func (c *Ctx) Watch(mask Level) CtxInterface {
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
	return c.Trace(LevelDebug, msg)
}

func (c *Ctx) Info(msg string) RecordInterface {
	return c.Trace(LevelInfo, msg)
}

func (c *Ctx) Warn(msg string) RecordInterface {
	return c.Trace(LevelWarn, msg)
}

func (c *Ctx) Error(msg string) RecordInterface {
	return c.Trace(LevelError, msg)
}

func (c *Ctx) Fatal(msg string) RecordInterface {
	return c.Trace(LevelFatal, msg)
}

func (c *Ctx) Assert(msg string) RecordInterface {
	return c.Trace(LevelAssert, msg)
}

func (c *Ctx) Trace(mask Level, msg string) RecordInterface {
	level := c.lmask & mask
	if level > 0 {
		return c.record(level, msg)
	}
	return DummyRecord{}
}

func (c *Ctx) DebugIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelDebug, msg)
}

func (c *Ctx) InfoIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelInfo, msg)
}

func (c *Ctx) WarnIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelWarn, msg)
}

func (c *Ctx) ErrorIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelError, msg)
}

func (c *Ctx) FatalIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelFatal, msg)
}

func (c *Ctx) AssertIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelAssert, msg)
}

func (c *Ctx) TraceIf(cond bool, mask Level, msg string) RecordInterface {
	if !cond {
		return DummyRecord{}
	}
	return c.Trace(mask, msg)
}

func (c *Ctx) record(level Level, msg string) *Record {
	r := c.newRecord(0)
	r.lp = c.log(level, "", msg, nil, false, EntryChapter, 0, r.id)
	return r
}

func (c *Ctx) dlog(level Level, name string, val interface{}, typ EntryType, tid, rid uint32) (dp uintptr) {
	c.lock()
	c.flushDL()
	dp = c.dlogLF(level, name, val, typ, tid, rid)
	c.unlock()
	return
}

func (c *Ctx) dlogLF(level Level, name string, val interface{}, typ EntryType, tid, rid uint32) (dp uintptr) {
	var tt int64
	if c.cl != nil {
		tt = c.cl.Now().UnixNano()
	} else {
		tt = time.Now().UnixNano()
	}
	e := c.dlogAcqLF()
	e.ll = level
	e.tp = typ
	e.tt = tt
	e.tid = tid
	e.rid = rid
	e.k = name
	e.v = val
	dp = uintptr(unsafe.Pointer(e))
	return
}

func (c *Ctx) dlogAcqLF() (e *dentry) {
	if c.dlbc < len(c.dlb) {
		e = &c.dlb[c.dlbc]
	} else {
		c.dlb = append(c.dlb, dentry{})
		e = &c.dlb[len(c.dlb)-1]
	}
	c.dlbc++
	return
}

func (c *Ctx) log(level Level, name string, val interface{}, m Marshaller, ind bool, typ EntryType, tid, rid uint32) (lp uintptr) {
	c.lock()
	c.flushDL()
	lp = c.logLF(level, name, val, m, ind, typ, tid, rid)
	c.unlock()
	return
}

func (c *Ctx) logLF(level Level, name string, val interface{}, m Marshaller, ind bool, typ EntryType, tid, rid uint32) (lp uintptr) {
	off := len(c.buf)
	var k Entry64
	if l := len(name); l > 0 {
		c.buf = append(c.buf, name...)
		k.Encode(uint32(off), uint32(off+l))
	}

	off = len(c.buf)
	var v Entry64
	c.bb.Reset()
	if typ == EntryLog {
		if vb, err := c.getm(m).Marshal(&c.bb, val, ind); err == nil {
			c.buf = append(c.buf, vb...)
			v.Encode(uint32(off), uint32(off+len(vb)))
		}
	} else {
		var err error
		if c.buf, err = x2bytes.ToBytes(c.buf, val); err == nil {
			v.Encode(uint32(off), uint32(len(c.buf)))
		}
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
	return uintptr(unsafe.Pointer(&c.lb[len(c.lb)-1]))
}

func (c *Ctx) flushDL() {
	for i := 0; i < c.dlbc; i++ {
		e := &c.dlb[i]
		var (
			m   Marshaller
			ind bool
		)
		for j := 0; j < len(e.opt); j++ {
			o := &e.opt[j]
			switch o.k {
			case OptionMarshaller:
				m, _ = o.v.(Marshaller)
			case OptionIndent:
				ind, _ = o.v.(bool)
			}
		}
		c.logLF(e.ll, e.k, e.v, m, ind, e.tp, e.tid, e.rid)
		e.reset()
	}
	c.dlbc = 0
}

func (c *Ctx) IsDummy() bool {
	return false
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

func (c *Ctx) AcquireThreadID(id uint32) ThreadInterface {
	t := c.newThread(0)
	t.SetID(id)
	c.log(LevelDebug, "", t.id, nil, false, EntryAcquireThread, 0, 0)
	return t
}

func (c *Ctx) ReleaseThread(thread ThreadInterface) CtxInterface {
	c.log(LevelDebug, "", thread.GetID(), nil, false, EntryReleaseThread, 0, 0)
	return c
}

func (c *Ctx) Reset() *Ctx {
	c.lmask = LevelAll
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
	for i := 0; i < c.dlbc; i++ {
		c.dlb[i].reset()
	}
	c.dlbc = 0
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
