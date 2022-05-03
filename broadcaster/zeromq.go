package broadcaster

import (
	"context"
	"sync"

	"github.com/pebbe/zmq4"
)

const (
	zmqDefaultTopic = "traceq"
	zmqDefaultHWM   = 1000
)

type ZeroMQ struct {
	Addr  string
	HWM   int
	Topic string

	once   sync.Once
	topic  []byte
	stream chan []byte
	err    error
}

func (b *ZeroMQ) Broadcast(ctx context.Context, p []byte) (n int, err error) {
	b.once.Do(func() {
		if len(b.Topic) == 0 {
			b.Topic = zmqDefaultTopic
		}
		b.topic = []byte(b.Topic)

		var (
			ztx *zmq4.Context
			zsk *zmq4.Socket
		)
		if ztx, b.err = zmq4.NewContext(); b.err != nil {
			return
		}
		if zsk, b.err = ztx.NewSocket(zmq4.PUB); b.err != nil {
			return
		}
		if b.HWM == 0 {
			b.HWM = zmqDefaultHWM
		}
		if b.err = zsk.SetSndhwm(b.HWM); b.err != nil {
			return
		}
		if b.err = zsk.Connect(b.Addr); b.err != nil {
			return
		}

		b.stream = make(chan []byte)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case p := <-b.stream:
					_, _ = zsk.SendBytes(b.topic, zmq4.SNDMORE)
					_, _ = zsk.SendBytes(p, 0)
					_ = p
				}
			}
		}()
	})

	if b.err != nil {
		err = b.err
		return
	}

	b.stream <- p

	return
}
