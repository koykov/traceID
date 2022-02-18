package traceID

import "time"

const (
	Version uint16 = 1
)

type CtxInterface interface {
	SetBroadcastTimeout(time.Duration) CtxInterface
	SetClock(Clock) CtxInterface
	SetMarshaller(Marshaller) CtxInterface
	SetID(string) CtxInterface
	Subject(string) CtxInterface
	Log(string, interface{}) CtxInterface
	LogWithOptions(string, interface{}, Options) CtxInterface
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) CtxInterface
	Flush() error
}

type ThreadInterface interface {
	GetID() uint32
	Subject(string) ThreadInterface
	Log(string, interface{}) ThreadInterface
	LogWithOptions(string, interface{}, Options) ThreadInterface
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
	Flush() error
}
