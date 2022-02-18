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

func (d DummyCtx) SetBroadcastTimeout(time.Duration) CtxInterface           { return d }
func (d DummyCtx) SetClock(Clock) CtxInterface                              { return d }
func (d DummyCtx) SetMarshaller(Marshaller) CtxInterface                    { return d }
func (d DummyCtx) SetID(string) CtxInterface                                { return d }
func (d DummyCtx) AcquireThread() ThreadInterface                           { return DummyThread{} }
func (d DummyCtx) ReleaseThread(ThreadInterface) CtxInterface               { return d }
func (d DummyCtx) Debug(string) CtxInterface                                { return d }
func (d DummyCtx) Info(string) CtxInterface                                 { return d }
func (d DummyCtx) Warn(string) CtxInterface                                 { return d }
func (d DummyCtx) Error(string) CtxInterface                                { return d }
func (d DummyCtx) Fatal(string) CtxInterface                                { return d }
func (d DummyCtx) Var(string, interface{}) CtxInterface                     { return d }
func (d DummyCtx) VarWithOptions(string, interface{}, Options) CtxInterface { return d }
func (d DummyCtx) Flush() error                                             { return nil }

type DummyThread struct{}

func (t DummyThread) GetID() uint32                                               { return 0 }
func (t DummyThread) Debug(string) ThreadInterface                                { return &t }
func (t DummyThread) Info(string) ThreadInterface                                 { return &t }
func (t DummyThread) Warn(string) ThreadInterface                                 { return &t }
func (t DummyThread) Error(string) ThreadInterface                                { return &t }
func (t DummyThread) Fatal(string) ThreadInterface                                { return &t }
func (t DummyThread) Var(string, interface{}) ThreadInterface                     { return &t }
func (t DummyThread) VarWithOptions(string, interface{}, Options) ThreadInterface { return &t }
func (t DummyThread) Flush() error                                                { return nil }
func (t DummyThread) AcquireThread() ThreadInterface                              { return DummyThread{} }
func (t DummyThread) ReleaseThread(ThreadInterface) ThreadInterface               { return t }

type DummyBroadcast struct{}

func (d DummyBroadcast) Broadcast(context.Context, []byte) (int, error) { return 0, nil }

type DummyListener struct{}

func (d DummyListener) SetAddr(string)                            {}
func (d DummyListener) Listen(context.Context, chan []byte) error { return nil }
