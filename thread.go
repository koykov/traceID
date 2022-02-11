package traceID

type Thread struct {
	id uint32
	cp uintptr
}

func (t Thread) Subject(subject string) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF("", subject, nil, EntrySubject)
	ctx.unlock()
	return &t
}

func (t Thread) Log(key string, val interface{}) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF(key, val, nil, EntryLog)
	ctx.unlock()
	return &t
}

func (t Thread) LogWM(key string, val interface{}, m Marshaller) ThreadInterface {
	ctx := t.indirectCtx()
	if ctx == nil {
		return &t
	}
	ctx.lock()
	ctx.logLF(key, val, m, EntryLog)
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

func (t Thread) indirectCtx() *Ctx {
	if t.cp == 0 {
		return nil
	}
	return indirectCtx(t.cp)
}
