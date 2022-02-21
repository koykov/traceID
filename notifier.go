package traceID

import "context"

type Notifier interface {
	SetAddr(string)
	Notify(context.Context, string) error
}
