package listener

import (
	"context"

	"github.com/koykov/traceID"
	"github.com/pebbe/zmq4"
)

type ZeroMQ struct {
	base
}

func (l ZeroMQ) Listen(ctx context.Context, out chan []byte) (err error) {
	if len(l.conf.Topic) == 0 {
		l.conf.Topic = traceID.DefaultZeroMQTopic
	}

	var (
		ztx *zmq4.Context
		zsk *zmq4.Socket
	)
	if ztx, err = zmq4.NewContext(); err != nil {
		return
	}
	if zsk, err = ztx.NewSocket(zmq4.SUB); err != nil {
		return
	}
	if l.conf.HWM == 0 {
		l.conf.HWM = traceID.DefaultZeroMQHWM
	}
	if err = zsk.SetSndhwm(int(l.conf.HWM)); err != nil {
		return
	}
	if err = zsk.Connect(l.conf.Addr); err != nil {
		return
	}
	if err = zsk.SetSubscribe(l.conf.Topic); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			_ = zsk.Close()
			return
		default:
			if _, err = zsk.RecvBytes(0); err != nil {
				continue
			}
			var p []byte
			if p, err = zsk.RecvBytes(0); err != nil || len(p) == 0 {
				continue
			}
			out <- p
		}
	}
}
