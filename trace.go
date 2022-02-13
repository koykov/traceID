package traceID

const (
	Version uint16 = 1
)

type CtxInterface interface {
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
	Subject(string) ThreadInterface
	Log(string, interface{}) ThreadInterface
	LogWM(string, interface{}, Marshaller) ThreadInterface
	Commit() error
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
}
