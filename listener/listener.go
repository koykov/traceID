package listener

import "github.com/koykov/traceID"

type listener struct {
	conf *traceID.ListenerConfig
}

func (l *listener) SetConfig(conf *traceID.ListenerConfig) {
	l.conf = conf
}
