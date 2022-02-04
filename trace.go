package traceID

const (
	Version uint16 = 1
)

type Interface interface {
	SetMarshaller(Marshaller) Interface
	SetID(string) Interface
	Subject(string) Interface
	Log(string, interface{}) Interface
	BeginTXN() Interface
	Commit() error
}
