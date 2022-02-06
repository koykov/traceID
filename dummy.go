package traceID

import "time"

type DummyClock struct{}

func (d DummyClock) Now() time.Time {
	return time.Unix(123, 456)
}

type DummyCtx struct{}

func (d DummyCtx) SetClock(Clock) Interface                        { return d }
func (d DummyCtx) SetMarshaller(Marshaller) Interface              { return d }
func (d DummyCtx) SetID(string) Interface                          { return d }
func (d DummyCtx) Subject(string) Interface                        { return d }
func (d DummyCtx) Log(string, interface{}) Interface               { return d }
func (d DummyCtx) LogWM(string, interface{}, Marshaller) Interface { return d }
func (d DummyCtx) BeginTXN() Interface                             { return d }
func (d DummyCtx) Commit() error                                   { return nil }

type DummyBroadcast struct{}

func (d DummyBroadcast) Broadcast([]byte) (int, error) { return 0, nil }
