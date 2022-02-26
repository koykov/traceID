package traceID

type Thread struct {
	ctxHeir
	id uint32
	rt uint32
}

func (t Thread) GetID() uint32 {
	return t.id
}

func (t Thread) Debug(msg string) RecordInterface {
	r := t.newRecord(LevelDebug, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (t Thread) Info(msg string) RecordInterface {
	r := t.newRecord(LevelInfo, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (t Thread) Warn(msg string) RecordInterface {
	r := t.newRecord(LevelWarn, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (t Thread) Error(msg string) RecordInterface {
	r := t.newRecord(LevelError, msg)
	if r == nil {
		return DummyRecord{}
	}
	return r
}

func (t Thread) Fatal(msg string) RecordInterface {
	r := t.newRecord(LevelFatal, msg)
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

func (t Thread) ReleaseThread(thread ThreadInterface) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.log(LevelDebug, "", thread.GetID(), nil, false, EntryReleaseThread, t.id, 0)
	return t
}

func (t Thread) newRecord(level LogLevel, msg string) *Record {
	ctx := t.indirectCtx()
	if ctx == nil {
		return nil
	}
	r := ctx.newRecord(t.id)
	ctx.log(level, "", msg, nil, false, EntryChapter, t.id, r.id)
	return r
}
