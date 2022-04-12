package broadcaster

import (
	"context"
)

type ZeroMQ struct {
	Addr string
}

func (b *ZeroMQ) Broadcast(ctx context.Context, p []byte) (n int, err error) {
	// todo implement me
	_, _ = ctx, p
	return
}
