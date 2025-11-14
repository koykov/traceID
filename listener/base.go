package listener

import "github.com/koykov/traceID"

type Base struct {
	conf *traceID.ListenerConfig
}

func (l *Base) SetConfig(conf *traceID.ListenerConfig) {
	l.conf = conf
}

func (l *Base) GetConfig() *traceID.ListenerConfig {
	return l.conf
}
