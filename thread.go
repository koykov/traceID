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
	ctx.log("", subject, nil, false, EntrySubject, t.id)
	return &t
}

func (t Thread) Log(key string, val interface{}) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(key, val, nil, false, EntryLog, t.id)
	return &t
}

func (t Thread) LogWithOptions(key string, val interface{}, opts Options) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(key, val, opts.Marshaller, opts.Indent, EntryLog, t.id)
	return &t
}

func (t Thread) Commit() (err error) {
	ctx := t.indirectCtx()
	if ctx == nil {
		err = ErrHomelessThread
		return
	}
	err = ctx.Commit()
	return
}

func (t Thread) AcquireThread() ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	id := atomic.AddUint32(&ctx.thc, 1)
	ctx.log("", id, nil, false, EntryAcquireThread, t.id)
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
	ctx.log("", thread.GetID(), nil, false, EntryReleaseThread, t.id)
	return t
}

func (t Thread) indirectCtx() *Ctx {
	if t.cp == 0 {
		return nil
	}
	return indirectCtx(t.cp)
}
