package traceID

import "time"

type LogLevel uint8

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal

	Version uint16 = 1
)

type CtxInterface interface {
	SetBroadcastTimeout(time.Duration) CtxInterface
	SetClock(Clock) CtxInterface
	SetMarshaller(Marshaller) CtxInterface
	SetService(string) CtxInterface
	SetID(string) CtxInterface
	Debug(string) RecordInterface
	Info(string) RecordInterface
	Warn(string) RecordInterface
	Error(string) RecordInterface
	Fatal(string) RecordInterface
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) CtxInterface
	Flush() error
}

type ThreadInterface interface {
	GetID() uint32
	Debug(string) RecordInterface
	Info(string) RecordInterface
	Warn(string) RecordInterface
	Error(string) RecordInterface
	Fatal(string) RecordInterface
	AcquireThread() ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
	Flush() error
}

type RecordInterface interface {
	Var(string, interface{}) RecordInterface
	VarWithOptions(string, interface{}, Options) RecordInterface
	VarIf(bool, string, interface{}) RecordInterface
	VarWithOptionsIf(bool, string, interface{}, Options) RecordInterface
	Flush() error
}
