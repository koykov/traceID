package traceID

import "github.com/koykov/indirect"

type ctxHeir struct {
	cp uintptr
}

func (h *ctxHeir) Flush() (err error) {
	ctx := h.indirectCtx()
	if ctx == nil {
		err = ErrHomelessThread
		return
	}
	err = ctx.Flush()
	return
}

func (h *ctxHeir) indirectCtx() *Ctx {
	if h.cp == 0 {
		return nil
	}
	up := indirect.ToUnsafePtr(h.cp)
	return (*Ctx)(up)
}
