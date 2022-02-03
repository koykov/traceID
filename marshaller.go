package traceID

import "io"

type Marshaller interface {
	Marshal(io.ReadWriter, interface{}) ([]byte, error)
}
