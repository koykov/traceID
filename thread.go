package traceID

type Thread struct {
	ctxHeir
	id uint32
	rt uint32
}

func (t *Thread) SetID(id uint32) ThreadInterface {
	t.id = id
	return t
}

func (t Thread) GetID() uint32 {
	return t.id
}

func (t Thread) Debug(msg string) RecordInterface {
	return t.Trace(LevelDebug, msg)
}

func (t Thread) Info(msg string) RecordInterface {
	return t.Trace(LevelInfo, msg)
}

func (t Thread) Warn(msg string) RecordInterface {
	return t.Trace(LevelWarn, msg)
}

func (t Thread) Error(msg string) RecordInterface {
	return t.Trace(LevelError, msg)
}

func (t Thread) Fatal(msg string) RecordInterface {
	return t.Trace(LevelFatal, msg)
}

func (t Thread) Assert(msg string) RecordInterface {
	return t.Trace(LevelAssert, msg)
}

func (t Thread) Trace(mask Level, msg string) RecordInterface {
	r := t.newRecord(mask, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (t Thread) DebugIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelDebug, msg)
}

func (t Thread) InfoIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelInfo, msg)
}

func (t Thread) WarnIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelWarn, msg)
}

func (t Thread) ErrorIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelError, msg)
}

func (t Thread) FatalIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelFatal, msg)
}

func (t Thread) AssertIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelAssert, msg)
}

func (t Thread) TraceIf(cond bool, mask Level, msg string) RecordInterface {
	if !cond {
		return DummyRecord{}
	}
	r := t.newRecord(mask, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (t Thread) AcquireThread() ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return DummyThread{}
	}
	return ctx.newThread(t.id)
}

func (t Thread) AcquireThreadID(id uint32) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return DummyThread{}
	}
	return ctx.newThread(t.id).SetID(id)
}

func (t Thread) ReleaseThread(thread ThreadInterface) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(LevelDebug, "", thread.GetID(), nil, false, EntryReleaseThread, t.id, ctx.nextRID())
	return &t
}

func (t Thread) newRecord(mask Level, msg string) *Record {
	ctx := t.indirectCtx()
	if ctx == nil {
		return nil
	}
	level := ctx.lmask & mask
	if level == 0 {
		return nil
	}
	r := ctx.newRecord(t.id)
	r.lp = ctx.log(level, "", msg, nil, false, EntryChapter, t.id, r.id)
	return r
}
