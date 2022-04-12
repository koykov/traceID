package broadcaster

import (
	"context"
)

type NATS struct {
	Addr string
}

func (b *NATS) Broadcast(ctx context.Context, p []byte) (n int, err error) {
	// todo implement me
	_, _ = ctx, p
	return
}
