package traceID

import "sync"

type CtxPool struct {
	p sync.Pool
}

var (
	CP CtxPool

	_, _, _ = AcquireCtx, ReleaseCtx, TryReleaseCtx
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

func AcquireCtx() *Ctx {
	return CP.Get()
}

func ReleaseCtx(ctx *Ctx) {
	CP.Put(ctx)
}

func TryReleaseCtx(ci CtxInterface) {
	if ctx, ok := interface{}(ci).(*Ctx); ok {
		ReleaseCtx(ctx)
	}
}
