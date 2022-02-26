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
	ctx.log(LevelDebug, name, val, nil, false, EntryLog, r.thid, r.id)
	return &r
}

func (r Record) VarWithOptions(name string, val interface{}, opts Options) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	ctx.log(LevelDebug, name, val, opts.Marshaller, opts.Indent, EntryLog, r.thid, r.id)
	return &r
}

func (r Record) VarIf(cond bool, name string, val interface{}) RecordInterface {
	if !cond {
		return &r
	}
	return r.Var(name, val)
}

func (r Record) VarWithOptionsIf(cond bool, name string, val interface{}, opts Options) RecordInterface {
	if !cond {
		return &r
	}
	return r.VarWithOptions(name, val, opts)
}
