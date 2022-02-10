package traceID

import "context"

type Listener interface {
	SetAddr(string)
	Listen(context.Context, chan []byte) error
}
