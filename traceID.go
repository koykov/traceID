package traceID

import "time"

type Level uint8
type Flag int

const (
	LevelDebug  Level = 1
	LevelInfo   Level = 1 << 1
	LevelWarn   Level = 1 << 2
	LevelError  Level = 1 << 3
	LevelFatal  Level = 1 << 4
	LevelAssert Level = 1 << 5
	LogAll            = LevelDebug | LevelInfo | LevelWarn | LevelError | LevelFatal | LevelAssert

	FlagOverwrite Flag = 0

	Version uint16 = 1
)

type CtxInterface interface {
	SetID(string) CtxInterface
	SetService(string) CtxInterface
	SetFlag(Flag, bool) CtxInterface
	Watch(Level) CtxInterface
	SetBroadcastTimeout(time.Duration) CtxInterface
	SetClock(Clock) CtxInterface
	SetMarshaller(Marshaller) CtxInterface
	SetLogger(Logger) CtxInterface
	Debug(string) RecordInterface
	Info(string) RecordInterface
	Warn(string) RecordInterface
	Error(string) RecordInterface
	Fatal(string) RecordInterface
	Assert(string) RecordInterface
	Trace(Level, string) RecordInterface
	AcquireThread() ThreadInterface
	AcquireThreadID(uint32) ThreadInterface
	ReleaseThread(ThreadInterface) CtxInterface
	Flush() error
}

type ThreadInterface interface {
	SetID(uint32) ThreadInterface
	GetID() uint32
	Debug(string) RecordInterface
	Info(string) RecordInterface
	Warn(string) RecordInterface
	Error(string) RecordInterface
	Fatal(string) RecordInterface
	Assert(string) RecordInterface
	Trace(Level, string) RecordInterface
	AcquireThread() ThreadInterface
	AcquireThreadID(uint32) ThreadInterface
	ReleaseThread(ThreadInterface) ThreadInterface
	Flush() error
}

type RecordInterface interface {
	Var(string, interface{}) RecordInterface
	VarIf(bool, string, interface{}) RecordInterface
	With(Option, interface{}) RecordInterface
	Err(error) RecordInterface
	ErrIf(bool, error) RecordInterface
	Flush() error
}

func (l Level) String() string {
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
	case LevelAssert:
		return "ASSERT"
	default:
		return "UNK"
	}
}

var _ = FlagOverwrite
