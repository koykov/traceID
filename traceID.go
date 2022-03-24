package traceID

import "time"

type LogLevel uint8
type Flag int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal

	FlagOverwrite Flag = 0

	Version uint16 = 1
)

type CtxInterface interface {
	SetFlag(Flag, bool) CtxInterface
	SetBroadcastTimeout(time.Duration) CtxInterface
	SetClock(Clock) CtxInterface
	SetMarshaller(Marshaller) CtxInterface
	SetLogger(Logger) CtxInterface
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
	Err(error) RecordInterface
	ErrIf(bool, error) RecordInterface
	Flush() error
}

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNK"
	}
}