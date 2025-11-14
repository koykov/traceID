package traceID

import (
	"context"
	"sync"
	"time"
)

type Broadcaster interface {
	SetConfig(config *BroadcasterConfig)
	GetConfig() *BroadcasterConfig
	Broadcast(context.Context, []byte) (int, error)
}

type BroadcasterConfig struct {
	Handler   string        `json:"handler"`
	Addr      string        `json:"addr"`
	Path      string        `json:"path,omitempty"`
	HWM       uint          `json:"hwm,omitempty"`
	Ping      uint          `json:"ping,omitempty"`
	PingDelay time.Duration `json:"pingDelay,omitempty"`
	Topic     string        `json:"topic,omitempty"`
}

var (
	bcs  []Broadcaster
	_, _ = RegisterBroadcaster, Broadcast
)

func RegisterBroadcaster(bc Broadcaster) {
	bcs = append(bcs, bc)
}

func Broadcast(message []byte) (err error) {
	return BroadcastWithTimeout(message, 0)
}

func BroadcastWithTimeout(message []byte, timeout time.Duration) (err error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < len(bcs); i++ {
		wg.Add(1)
		go func(ctx context.Context, bc Broadcaster) {
			defer wg.Done()
			_, _ = bc.Broadcast(ctx, message)
		}(ctx, bcs[i])
	}
	wg.Wait()

	return
}
