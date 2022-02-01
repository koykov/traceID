package traceID

type Interface interface {
	BeginTXN() Interface
	Add(interface{}) Interface
	AddKV(string, interface{}) Interface
	Commit() error
	CommitTo(Broadcaster) error
}
