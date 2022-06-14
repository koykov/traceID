package traceID

import "context"

type Listener interface {
	SetConfig(*ListenerConfig)
	GetConfig() *ListenerConfig
	Listen(context.Context, chan []byte) error
}

type ListenerConfig struct {
	Handler string `json:"handler"`
	Addr    string `json:"addr"`
	Path    string `json:"path,omitempty"`
	HWM     uint   `json:"hwm,omitempty"`
	Topic   string `json:"topic,omitempty"`
}
