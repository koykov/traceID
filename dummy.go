package traceID

import (
	"context"
	"time"
)

type DummyClock struct{}

func (d DummyClock) Now() time.Time {
	return time.Unix(123, 456)
}

type DummyCtx struct {
	DummyAssert
}

func (d DummyCtx) SetID(string) CtxInterface                       { return d }
func (d DummyCtx) GetID() string                                   { return "" }
func (d DummyCtx) SetService(string) CtxInterface                  { return d }
func (d DummyCtx) SetServiceWithStage(string, string) CtxInterface { return d }
func (d DummyCtx) SetStage(string) CtxInterface                    { return d }
func (d DummyCtx) SetFlag(Flag, bool) CtxInterface                 { return d }
func (d DummyCtx) Watch(Level) CtxInterface                        { return d }
func (d DummyCtx) SetBroadcastTimeout(time.Duration) CtxInterface  { return d }
func (d DummyCtx) SetClock(Clock) CtxInterface                     { return d }
func (d DummyCtx) SetMarshaller(Marshaller) CtxInterface           { return d }
func (d DummyCtx) SetLogger(Logger) CtxInterface                   { return d }
func (d DummyCtx) AcquireThread() ThreadInterface                  { return DummyThread{} }
func (d DummyCtx) AcquireThreadID(uint32) ThreadInterface          { return DummyThread{} }
func (d DummyCtx) ReleaseThread(ThreadInterface) CtxInterface      { return d }
func (d DummyCtx) Debug(string) RecordInterface                    { return DummyRecord{} }
func (d DummyCtx) Info(string) RecordInterface                     { return DummyRecord{} }
func (d DummyCtx) Warn(string) RecordInterface                     { return DummyRecord{} }
func (d DummyCtx) Error(string) RecordInterface                    { return DummyRecord{} }
func (d DummyCtx) Fatal(string) RecordInterface                    { return DummyRecord{} }
func (d DummyCtx) Assert(string) RecordInterface                   { return DummyRecord{} }
func (d DummyCtx) Trace(Level, string) RecordInterface             { return DummyRecord{} }
func (d DummyCtx) DebugIf(bool, string) RecordInterface            { return DummyRecord{} }
func (d DummyCtx) InfoIf(bool, string) RecordInterface             { return DummyRecord{} }
func (d DummyCtx) WarnIf(bool, string) RecordInterface             { return DummyRecord{} }
func (d DummyCtx) ErrorIf(bool, string) RecordInterface            { return DummyRecord{} }
func (d DummyCtx) FatalIf(bool, string) RecordInterface            { return DummyRecord{} }
func (d DummyCtx) AssertIf(bool, string) RecordInterface           { return DummyRecord{} }
func (d DummyCtx) TraceIf(bool, Level, string) RecordInterface     { return DummyRecord{} }
func (d DummyCtx) IsDummy() bool                                   { return true }
func (d DummyCtx) Flush() error                                    { return nil }

type DummyThread struct {
	DummyAssert
}

func (t DummyThread) SetID(uint32) ThreadInterface                  { return t }
func (t DummyThread) GetID() uint32                                 { return 0 }
func (t DummyThread) Debug(string) RecordInterface                  { return DummyRecord{} }
func (t DummyThread) Info(string) RecordInterface                   { return DummyRecord{} }
func (t DummyThread) Warn(string) RecordInterface                   { return DummyRecord{} }
func (t DummyThread) Error(string) RecordInterface                  { return DummyRecord{} }
func (t DummyThread) Fatal(string) RecordInterface                  { return DummyRecord{} }
func (t DummyThread) Assert(string) RecordInterface                 { return DummyRecord{} }
func (t DummyThread) Trace(Level, string) RecordInterface           { return DummyRecord{} }
func (t DummyThread) DebugIf(bool, string) RecordInterface          { return DummyRecord{} }
func (t DummyThread) InfoIf(bool, string) RecordInterface           { return DummyRecord{} }
func (t DummyThread) WarnIf(bool, string) RecordInterface           { return DummyRecord{} }
func (t DummyThread) ErrorIf(bool, string) RecordInterface          { return DummyRecord{} }
func (t DummyThread) FatalIf(bool, string) RecordInterface          { return DummyRecord{} }
func (t DummyThread) AssertIf(bool, string) RecordInterface         { return DummyRecord{} }
func (t DummyThread) TraceIf(bool, Level, string) RecordInterface   { return DummyRecord{} }
func (t DummyThread) Flush() error                                  { return nil }
func (t DummyThread) AcquireThread() ThreadInterface                { return DummyThread{} }
func (t DummyThread) AcquireThreadID(uint32) ThreadInterface        { return DummyThread{} }
func (t DummyThread) ReleaseThread(ThreadInterface) ThreadInterface { return t }

type DummyRecord struct{}

func (r DummyRecord) Slug(string) RecordInterface                     { return r }
func (r DummyRecord) Var(string, interface{}) RecordInterface         { return r }
func (r DummyRecord) VarIf(bool, string, interface{}) RecordInterface { return r }
func (r DummyRecord) With(Option, interface{}) RecordInterface        { return r }
func (r DummyRecord) Err(error) RecordInterface                       { return r }
func (r DummyRecord) ErrIf(bool, error) RecordInterface               { return r }
func (r DummyRecord) Comment(string) RecordInterface                  { return r }
func (r DummyRecord) CommentIf(bool, string) RecordInterface          { return r }
func (r DummyRecord) Flush() error                                    { return nil }

type DummyBroadcast struct{}

func (d DummyBroadcast) SetConfig(*BroadcasterConfig)                   {}
func (d DummyBroadcast) GetConfig() *BroadcasterConfig                  { return nil }
func (d DummyBroadcast) Broadcast(context.Context, []byte) (int, error) { return 0, nil }

type DummyListener struct{}

func (d DummyListener) SetConfig(*ListenerConfig)                 {}
func (d DummyListener) GetConfig() *ListenerConfig                { return nil }
func (d DummyListener) Listen(context.Context, chan []byte) error { return nil }
