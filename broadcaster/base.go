package broadcaster

import "github.com/koykov/traceID"

type Base struct {
	conf *traceID.BroadcasterConfig
}

func (b *Base) SetConfig(conf *traceID.BroadcasterConfig) {
	b.conf = conf
}

func (b *Base) GetConfig() *traceID.BroadcasterConfig {
	return b.conf
}
