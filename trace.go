package traceID

type Interface interface {
	BeginTXN() Interface
	Log(string, interface{}) Interface
	Commit() error
	CommitTo(Broadcaster) error
}
