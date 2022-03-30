package traceID

import (
	"context"
	"time"
)

type DummyClock struct{}

func (d DummyClock) Now() time.Time {
	return time.Unix(123, 456)
}

type DummyCtx struct{}

func (d DummyCtx) SetID(string) CtxInterface                      { return d }
func (d DummyCtx) SetService(string) CtxInterface                 { return d }
func (d DummyCtx) SetFlag(Flag, bool) CtxInterface                { return d }
func (d DummyCtx) Watch(LogLevel) CtxInterface                    { return d }
func (d DummyCtx) SetBroadcastTimeout(time.Duration) CtxInterface { return d }
func (d DummyCtx) SetClock(Clock) CtxInterface                    { return d }
func (d DummyCtx) SetMarshaller(Marshaller) CtxInterface          { return d }
func (d DummyCtx) SetLogger(Logger) CtxInterface                  { return d }
func (d DummyCtx) AcquireThread() ThreadInterface                 { return DummyThread{} }
func (d DummyCtx) ReleaseThread(ThreadInterface) CtxInterface     { return d }
func (d DummyCtx) Debug(string) RecordInterface                   { return DummyRecord{} }
func (d DummyCtx) Info(string) RecordInterface                    { return DummyRecord{} }
func (d DummyCtx) Warn(string) RecordInterface                    { return DummyRecord{} }
func (d DummyCtx) Error(string) RecordInterface                   { return DummyRecord{} }
func (d DummyCtx) Fatal(string) RecordInterface                   { return DummyRecord{} }
func (d DummyCtx) Assert(string) RecordInterface                  { return DummyRecord{} }
func (d DummyCtx) Log(LogLevel, string) RecordInterface           { return DummyRecord{} }
func (d DummyCtx) Flush() error                                   { return nil }

type DummyThread struct{}

func (t DummyThread) GetID() uint32                                 { return 0 }
func (t DummyThread) Debug(string) RecordInterface                  { return DummyRecord{} }
func (t DummyThread) Info(string) RecordInterface                   { return DummyRecord{} }
func (t DummyThread) Warn(string) RecordInterface                   { return DummyRecord{} }
func (t DummyThread) Error(string) RecordInterface                  { return DummyRecord{} }
func (t DummyThread) Fatal(string) RecordInterface                  { return DummyRecord{} }
func (t DummyThread) Flush() error                                  { return nil }
func (t DummyThread) AcquireThread() ThreadInterface                { return DummyThread{} }
func (t DummyThread) ReleaseThread(ThreadInterface) ThreadInterface { return t }

type DummyRecord struct{}

func (r DummyRecord) Var(string, interface{}) RecordInterface                             { return &r }
func (r DummyRecord) VarWithOptions(string, interface{}, Options) RecordInterface         { return &r }
func (r DummyRecord) VarIf(bool, string, interface{}) RecordInterface                     { return &r }
func (r DummyRecord) VarWithOptionsIf(bool, string, interface{}, Options) RecordInterface { return &r }
func (r DummyRecord) Err(error) RecordInterface                                           { return &r }
func (r DummyRecord) ErrIf(bool, error) RecordInterface                                   { return &r }
func (r DummyRecord) Flush() error                                                        { return nil }

type DummyBroadcast struct{}

func (d DummyBroadcast) Broadcast(context.Context, []byte) (int, error) { return 0, nil }

type DummyListener struct{}

func (d DummyListener) SetConfig(*ListenerConfig)                 {}
func (d DummyListener) Listen(context.Context, chan []byte) error { return nil }
