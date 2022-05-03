package traceID

import "context"

type Listener interface {
	SetConfig(*ListenerConfig)
	Listen(context.Context, chan []byte) error
}

type ListenerConfig struct {
	Handler string `json:"handler"`
	Addr    string `json:"addr"`
	Path    string `json:"path"`
	BufSize uint   `json:"buf_size"`
}
