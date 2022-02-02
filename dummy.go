package traceID

type DummyCtx struct{}

func (d DummyCtx) BeginTXN() Interface               { return d }
func (d DummyCtx) Log(string, interface{}) Interface { return d }
func (d DummyCtx) Commit() error                     { return nil }
func (d DummyCtx) CommitTo(Broadcaster) error        { return nil }

type DummyBroadcast struct{}

func (d DummyBroadcast) Broadcast([]byte) (int, error) { return 0, nil }
