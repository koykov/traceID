package traceID

type Thread struct {
	id  uint32
	ctx *Ctx
}

func (t Thread) Subject(subject string) ThreadInterface {
	if t.ctx == nil {
		return &t
	}
	t.ctx.lock()
	t.ctx.logLF("", subject, nil, EntrySubject)
	t.ctx.unlock()
	return &t
}

func (t Thread) Log(key string, val interface{}) ThreadInterface {
	if t.ctx == nil {
		return &t
	}
	t.ctx.lock()
	t.ctx.logLF(key, val, nil, EntryLog)
	t.ctx.unlock()
	return &t
}

func (t Thread) LogWM(key string, val interface{}, m Marshaller) ThreadInterface {
	if t.ctx == nil {
		return &t
	}
	t.ctx.lock()
	t.ctx.logLF(key, val, m, EntryLog)
	t.ctx.unlock()
	return &t
}

func (t Thread) Commit() (err error) {
	if t.ctx == nil {
		err = ErrHomelessThread
		return
	}
	t.ctx.lock()
	err = t.ctx.Commit()
	t.ctx.unlock()
	return
}
