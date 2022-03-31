package traceID

type Record struct {
	ctxHeir
	id   uint32
	thid uint32
}

func (r Record) Var(name string, val interface{}) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	ctx.dlog(LevelDebug, name, val, EntryLog, r.thid, r.id)
	return &r
}

func (r Record) VarIf(cond bool, name string, val interface{}) RecordInterface {
	if !cond {
		return &r
	}
	return r.Var(name, val)
}

func (r Record) With(name Option, value interface{}) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	ctx.addOpt(r.id, name, value)
	return &r
}

func (r Record) Err(err error) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	ctx.dlog(LevelDebug, "", err.Error(), EntryLog, r.thid, r.id)
	return &r
}

func (r Record) ErrIf(cond bool, err error) RecordInterface {
	if !cond {
		return r
	}
	return r.Err(err)
}
