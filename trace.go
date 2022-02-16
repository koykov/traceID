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
	Commit() error
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) CtxInterface
}

type ThreadInterface interface {
	GetID() uint32
	Subject(string) ThreadInterface
	Log(string, interface{}) ThreadInterface
	LogWithOptions(string, interface{}, Options) ThreadInterface
	Commit() error
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
}
