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

	once  sync.Once
	ctx   *zmq4.Context
	sock  *zmq4.Socket
	topic []byte
	err   error
}

func (b *ZeroMQ) Broadcast(_ context.Context, p []byte) (n int, err error) {
	b.once.Do(func() {
		if len(b.Topic) == 0 {
			b.Topic = zmqDefaultTopic
		}
		b.topic = []byte(b.Topic)

		if b.ctx, b.err = zmq4.NewContext(); b.err != nil {
			return
		}
		if b.sock, b.err = b.ctx.NewSocket(zmq4.PUB); b.err != nil {
			return
		}
		if b.HWM == 0 {
			b.HWM = zmqDefaultHWM
		}
		if b.err = b.sock.SetSndhwm(b.HWM); b.err != nil {
			return
		}
		if b.err = b.sock.Connect(b.Addr); b.err != nil {
			return
		}
	})

	if b.err != nil {
		err = b.err
		return
	}

	if n, err = b.sock.SendBytes(b.topic, zmq4.SNDMORE); err != nil {
		return
	}
	var n1 int
	if n1, err = b.sock.SendBytes(p, 0); err != nil {
		return
	}
	n += n1

	return
}
