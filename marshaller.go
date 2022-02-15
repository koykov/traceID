package traceID

import (
	"fmt"
	"io"
)

type Marshaller interface {
	Marshal(io.ReadWriter, interface{}, bool) ([]byte, error)
}

var (
	defaultMarshaller = &mfmt{}
)

type mfmt struct{}

func (m mfmt) Marshal(rw io.ReadWriter, x interface{}, _ bool) ([]byte, error) {
	if _, err := fmt.Fprint(rw, x); err != nil {
		return nil, err
	}
	return io.ReadAll(rw)
}
