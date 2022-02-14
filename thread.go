package traceID

import (
	"sync/atomic"
)

type Thread struct {
	id uint32
	rt uint32
	cp uintptr
}

func (t Thread) GetID() uint32 {
	return t.id
}

func (t Thread) Subject(subject string) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF("", subject, nil, EntrySubject, t.id)
	ctx.unlock()
	return &t
}

func (t Thread) Log(key string, val interface{}) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF(key, val, nil, EntryLog, t.id)
	ctx.unlock()
	return &t
}

func (t Thread) LogWM(key string, val interface{}, m Marshaller) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF(key, val, m, EntryLog, t.id)
	ctx.unlock()
	return &t
}

func (t Thread) Commit() (err error) {
	ctx := t.indirectCtx()
	if ctx == nil {
		err = ErrHomelessThread
		return
	}
	ctx.lock()
	err = ctx.Commit()
	ctx.unlock()
	return
}

func (t Thread) AcquireThread() ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	id := atomic.AddUint32(&ctx.thc, 1)
	ctx.lock()
	ctx.logLF("", id, nil, EntryAcquireThread, t.id)
	ctx.unlock()
	return &Thread{
		id: id,
		rt: t.id,
		cp: t.cp,
	}
}

func (t Thread) ReleaseThread(thread ThreadInterface) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF("", thread.GetID(), nil, EntryReleaseThread, t.id)
	ctx.unlock()
	return t
}

func (t Thread) indirectCtx() *Ctx {
	if t.cp == 0 {
		return nil
	}
	return indirectCtx(t.cp)
}
