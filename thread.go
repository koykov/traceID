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

func (t Thread) Debug(message string) ThreadInterface {
	t.chapter(LevelDebug, message)
	return &t
}

func (t Thread) Info(message string) ThreadInterface {
	t.chapter(LevelInfo, message)
	return &t
}

func (t Thread) Warn(message string) ThreadInterface {
	t.chapter(LevelWarn, message)
	return &t
}

func (t Thread) Error(message string) ThreadInterface {
	t.chapter(LevelError, message)
	return &t
}

func (t Thread) Fatal(message string) ThreadInterface {
	t.chapter(LevelFatal, message)
	return &t
}

func (t Thread) Var(key string, val interface{}) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(LevelDebug, key, val, nil, false, EntryLog, t.id)
	return &t
}

func (t Thread) VarWithOptions(key string, val interface{}, opts Options) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(LevelDebug, key, val, opts.Marshaller, opts.Indent, EntryLog, t.id)
	return &t
}

func (t Thread) VarIf(condition bool, key string, val interface{}) ThreadInterface {
	if !condition {
		return &t
	}
	return t.Var(key, val)
}

func (t Thread) VarWithOptionsIf(condition bool, key string, val interface{}, opts Options) ThreadInterface {
	if !condition {
		return &t
	}
	return t.VarWithOptions(key, val, opts)
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

func (t Thread) chapter(level LogLevel, message string) {
	ctx := t.indirectCtx()
	if ctx == nil {
		return
	}
	ctx.log(level, "", message, nil, false, EntryChapter, t.id)
}

func (t Thread) indirectCtx() *Ctx {
	if t.cp == 0 {
		return nil
	}
	return indirectCtx(t.cp)
}
