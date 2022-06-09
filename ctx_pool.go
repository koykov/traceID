package traceID

import "sync"

type CtxPool struct {
	p sync.Pool
}

var (
	CP CtxPool

	_, _, _ = AcquireCtx, ReleaseCtx, AcquireCtxWithID
)

func (p *CtxPool) Get() *Ctx {
	v := p.p.Get()
	if v != nil {
		if c, ok := v.(*Ctx); ok {
			return c
		}
	}
	return NewCtx()
}

func (p *CtxPool) Put(ctx *Ctx) {
	ctx.Reset()
	p.p.Put(ctx)
}

func AcquireCtx() CtxInterface {
	return CP.Get()
}

func AcquireCtxWithID(id string) CtxInterface {
	if len(id) == 0 {
		return DummyCtx{}
	}
	ctx := AcquireCtx()
	ctx.SetID(id)
	return ctx
}

func ReleaseCtx(ctx CtxInterface) {
	tryReleaseCtx(ctx)
}

func tryReleaseCtx(ci CtxInterface) {
	if ctx, ok := interface{}(ci).(*Ctx); ok {
		CP.Put(ctx)
	}
}
