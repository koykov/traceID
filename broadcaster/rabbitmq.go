package broadcaster

import "context"

type RabbitMQ struct {
	Addr string
}

func (b *RabbitMQ) Broadcast(ctx context.Context, p []byte) (n int, err error) {
	// todo implement me
	_, _ = ctx, p
	return
}
