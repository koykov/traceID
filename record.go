package traceID

import "github.com/koykov/indirect"

type Record struct {
	ctxHeir
	id   uint32
	lp   uintptr
	dp   uintptr
	thid uint32
}

func (r Record) Slug(slug string) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	if r.lp == 0 {
		return &r
	}
	e := (*entry)(indirect.ToUnsafePtr(r.lp))
	lo := len(ctx.buf)
	ctx.buf = append(ctx.buf, slug...)
	hi := len(ctx.buf)
	e.k.Encode(uint32(lo), uint32(hi))
	return &r
}

func (r Record) Var(name string, val interface{}) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	r.dp = ctx.dlog(LevelDebug, name, val, EntryLog, r.thid, r.id)
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
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	if r.dp == 0 {
		return &r
	}
	e := (*dentry)(indirect.ToUnsafePtr(r.dp))
	e.opt = append(e.opt, optionKV{k: name, v: value})
	return &r
}

func (r Record) Err(err error) RecordInterface {
	ctx := r.indirectCtx()
	if ctx == nil {
		return &r
	}
	r.dp = ctx.dlog(LevelDebug, "", err.Error(), EntryLog, r.thid, r.id)
	return &r
}

func (r Record) ErrIf(cond bool, err error) RecordInterface {
	if !cond {
		return r
	}
	return r.Err(err)
}
