package listener

import "github.com/koykov/traceID"

type base struct {
	conf *traceID.ListenerConfig
}

func (l *base) SetConfig(conf *traceID.ListenerConfig) {
	l.conf = conf
}
