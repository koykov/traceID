package listener

import (
	"context"

	"github.com/koykov/traceID"
	"github.com/pebbe/zmq4"
)

const (
	zmqDefaultTopic = "traceq"
	zmqDefaultHWM   = 1000
)

type ZeroMQ struct {
	listener
	Addr  string
	HWM   int
	Topic string
}

func (l *ZeroMQ) SetConfig(conf *traceID.ListenerConfig) {
	l.listener.SetConfig(conf)
	l.HWM = int(conf.BufSize)
	l.Topic = conf.Path
}

func (l ZeroMQ) Listen(ctx context.Context, out chan []byte) (err error) {
	if len(l.Topic) == 0 {
		l.Topic = zmqDefaultTopic
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
	if l.HWM == 0 {
		l.HWM = zmqDefaultHWM
	}
	if err = zsk.SetSndhwm(l.HWM); err != nil {
		return
	}
	if err = zsk.Connect(l.Addr); err != nil {
		return
	}
	if err = zsk.SetSubscribe(l.Topic); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
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
