package traceID

import "time"

type LogLevel uint8

const (
	Version uint16 = 1

	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type CtxInterface interface {
	SetBroadcastTimeout(time.Duration) CtxInterface
	SetClock(Clock) CtxInterface
	SetMarshaller(Marshaller) CtxInterface
	SetService(string) CtxInterface
	SetID(string) CtxInterface
	Debug(string) CtxInterface
	Info(string) CtxInterface
	Warn(string) CtxInterface
	Error(string) CtxInterface
	Fatal(string) CtxInterface
	Var(string, interface{}) CtxInterface
	VarWithOptions(string, interface{}, Options) CtxInterface
	VarIf(bool, string, interface{}) CtxInterface
	VarWithOptionsIf(bool, string, interface{}, Options) CtxInterface
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) CtxInterface
	Flush() error
}

type ThreadInterface interface {
	GetID() uint32
	Debug(string) ThreadInterface
	Info(string) ThreadInterface
	Warn(string) ThreadInterface
	Error(string) ThreadInterface
	Fatal(string) ThreadInterface
	Var(string, interface{}) ThreadInterface
	VarWithOptions(string, interface{}, Options) ThreadInterface
	VarIf(bool, string, interface{}) ThreadInterface
	VarWithOptionsIf(bool, string, interface{}, Options) ThreadInterface
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
	Flush() error
}
