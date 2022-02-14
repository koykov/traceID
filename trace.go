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
	LogWM(string, interface{}, Marshaller) CtxInterface
	Commit() error
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) CtxInterface
}

type ThreadInterface interface {
	GetID() uint32
	Subject(string) ThreadInterface
	Log(string, interface{}) ThreadInterface
	LogWM(string, interface{}, Marshaller) ThreadInterface
	Commit() error
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
}
