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
	hwm   int
	topic string
}

func (l *ZeroMQ) SetConfig(conf *traceID.ListenerConfig) {
	l.listener.SetConfig(conf)
	l.hwm = int(conf.BufSize)
	l.topic = conf.Path
}

func (l ZeroMQ) Listen(ctx context.Context, out chan []byte) (err error) {
	if len(l.topic) == 0 {
		l.topic = zmqDefaultTopic
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
	if l.hwm == 0 {
		l.hwm = zmqDefaultHWM
	}
	if err = zsk.SetSndhwm(l.hwm); err != nil {
		return
	}
	if err = zsk.Connect(l.conf.Addr); err != nil {
		return
	}
	if err = zsk.SetSubscribe(l.topic); err != nil {
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
