package traceID

import "time"

type LogLevel uint8
type Flag int

const (
	LevelDebug  LogLevel = 1
	LevelInfo   LogLevel = 1 << 1
	LevelWarn   LogLevel = 1 << 2
	LevelError  LogLevel = 1 << 3
	LevelFatal  LogLevel = 1 << 4
	LevelAssert LogLevel = 1 << 5
	LogAll               = LevelDebug | LevelInfo | LevelWarn | LevelError | LevelFatal | LevelAssert

	FlagOverwrite Flag = 0

	Version uint16 = 1
)

type CtxInterface interface {
	SetID(string) CtxInterface
	SetService(string) CtxInterface
	SetFlag(Flag, bool) CtxInterface
	Watch(LogLevel) CtxInterface
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
	Log(LogLevel, string) RecordInterface
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

var _ = FlagOverwrite
