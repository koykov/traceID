package traceID

type Listener interface {
	Listen(string, chan []byte) error
}
