package traceID

import (
	"io"

	"github.com/koykov/traceID/marshaller"
)

type Marshaller interface {
	Marshal(io.ReadWriter, interface{}) ([]byte, error)
}

var (
	defaultMarshaller = &marshaller.JSON{}
)
