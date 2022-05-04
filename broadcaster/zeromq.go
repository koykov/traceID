package broadcaster

import (
	"context"
	"sync"

	"github.com/koykov/traceID"
	"github.com/pebbe/zmq4"
)

type ZeroMQ struct {
	base
	once  sync.Once
	ctx   *zmq4.Context
	sock  *zmq4.Socket
	topic []byte
	err   error
}

func (b *ZeroMQ) Broadcast(_ context.Context, p []byte) (n int, err error) {
	b.once.Do(func() {
		if len(b.conf.Topic) == 0 {
			b.conf.Topic = traceID.DefaultZeroMQTopic
		}
		b.topic = []byte(b.conf.Topic)

		if b.ctx, b.err = zmq4.NewContext(); b.err != nil {
			return
		}
		if b.sock, b.err = b.ctx.NewSocket(zmq4.PUB); b.err != nil {
			return
		}
		if b.conf.HWM == 0 {
			b.conf.HWM = traceID.DefaultZeroMQHWM
		}
		if b.err = b.sock.SetSndhwm(int(b.conf.HWM)); b.err != nil {
			return
		}
		if b.err = b.sock.Connect(b.conf.Addr); b.err != nil {
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
