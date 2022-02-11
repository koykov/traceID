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

func (d DummyCtx) SetClock(Clock) CtxInterface                        { return d }
func (d DummyCtx) SetMarshaller(Marshaller) CtxInterface              { return d }
func (d DummyCtx) SetID(string) CtxInterface                          { return d }
func (d DummyCtx) Thread() ThreadInterface                            { return DummyThread{} }
func (d DummyCtx) Subject(string) CtxInterface                        { return d }
func (d DummyCtx) Log(string, interface{}) CtxInterface               { return d }
func (d DummyCtx) LogWM(string, interface{}, Marshaller) CtxInterface { return d }
func (d DummyCtx) Commit() error                                      { return nil }

type DummyThread struct{}

func (t DummyThread) Subject(string) ThreadInterface                        { return &t }
func (t DummyThread) Log(string, interface{}) ThreadInterface               { return &t }
func (t DummyThread) LogWM(string, interface{}, Marshaller) ThreadInterface { return &t }
func (t DummyThread) Commit() error                                         { return nil }

type DummyBroadcast struct{}

func (d DummyBroadcast) Broadcast([]byte) (int, error) { return 0, nil }

type DummyListener struct{}

func (d DummyListener) SetAddr(string)                            {}
func (d DummyListener) Listen(context.Context, chan []byte) error { return nil }
