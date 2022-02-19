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

func (t Thread) Debug(msg string) ThreadInterface {
	t.chapter(LevelDebug, msg)
	return &t
}

func (t Thread) Info(msg string) ThreadInterface {
	t.chapter(LevelInfo, msg)
	return &t
}

func (t Thread) Warn(msg string) ThreadInterface {
	t.chapter(LevelWarn, msg)
	return &t
}

func (t Thread) Error(msg string) ThreadInterface {
	t.chapter(LevelError, msg)
	return &t
}

func (t Thread) Fatal(msg string) ThreadInterface {
	t.chapter(LevelFatal, msg)
	return &t
}

func (t Thread) Var(name string, val interface{}) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(LevelDebug, name, val, nil, false, EntryLog, t.id)
	return &t
}

func (t Thread) VarWithOptions(name string, val interface{}, opts Options) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(LevelDebug, name, val, opts.Marshaller, opts.Indent, EntryLog, t.id)
	return &t
}

func (t Thread) VarIf(cond bool, name string, val interface{}) ThreadInterface {
	if !cond {
		return &t
	}
	return t.Var(name, val)
}

func (t Thread) VarWithOptionsIf(cond bool, name string, val interface{}, opts Options) ThreadInterface {
	if !cond {
		return &t
	}
	return t.VarWithOptions(name, val, opts)
}

func (t Thread) Flush() (err error) {
	ctx := t.indirectCtx()
	if ctx == nil {
		err = ErrHomelessThread
		return
	}
	err = ctx.Flush()
	return
}

func (t Thread) AcquireThread() ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return DummyThread{}
	}
	id := atomic.AddUint32(&ctx.thc, 1)
	ctx.log(LevelDebug, "", id, nil, false, EntryAcquireThread, t.id)
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
	ctx.log(LevelDebug, "", thread.GetID(), nil, false, EntryReleaseThread, t.id)
	return t
}

func (t Thread) chapter(level LogLevel, msg string) {
	ctx := t.indirectCtx()
	if ctx == nil {
		return
	}
	ctx.log(level, "", msg, nil, false, EntryChapter, t.id)
}

func (t Thread) indirectCtx() *Ctx {
	if t.cp == 0 {
		return nil
	}
	return indirectCtx(t.cp)
}
