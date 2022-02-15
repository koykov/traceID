package traceID

import (
	"fmt"
	"io"
)

type Marshaller interface {
	Marshal(io.ReadWriter, interface{}) ([]byte, error)
	MarshalIndent(io.ReadWriter, interface{}) ([]byte, error)
}

var (
	defaultMarshaller = &mfmt{}
)

type mfmt struct{}

func (m mfmt) Marshal(rw io.ReadWriter, x interface{}) ([]byte, error) {
	if _, err := fmt.Fprint(rw, x); err != nil {
		return nil, err
	}
	return io.ReadAll(rw)
}

func (m mfmt) MarshalIndent(rw io.ReadWriter, x interface{}) ([]byte, error) {
	return m.Marshal(rw, x)
}
