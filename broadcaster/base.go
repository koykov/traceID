package broadcaster

import "github.com/koykov/traceID"

type base struct {
	conf *traceID.BroadcasterConfig
}

func (b *base) SetConfig(conf *traceID.BroadcasterConfig) {
	b.conf = conf
}
