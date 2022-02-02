package traceID

type Interface interface {
	SetComponent(string) Interface
	SetID(string) Interface
	Subject(string) Interface
	Log(string, interface{}) Interface
	Push() error
	BeginTXN() Interface
	Commit() error
}
