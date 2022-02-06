package traceID

const (
	Version uint16 = 1
)

type Interface interface {
	SetClock(Clock) Interface
	SetMarshaller(Marshaller) Interface
	SetID(string) Interface
	Subject(string) Interface
	Log(string, interface{}) Interface
	LogWM(string, interface{}, Marshaller) Interface
	BeginTXN() Interface
	Commit() error
}
