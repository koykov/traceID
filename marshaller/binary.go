package marshaller

import (
	"fmt"
	"io"

	"github.com/koykov/bytealg"
	"github.com/koykov/x2bytes"
)

type Binary struct{}

func (m Binary) Marshal(rw io.ReadWriter, x interface{}, _ bool) (b []byte, err error) {
	var p []byte
	if p, err = x2bytes.ToBytes(p, x); err != nil || len(p) == 0 {
		return
	}
	for i := 0; i < len(p); i++ {
		_, _ = fmt.Fprintf(rw, "%02X", p[i])
		_, _ = fmt.Fprintf(rw, " ")
	}
	b, err = io.ReadAll(rw)
	b = bytealg.TrimRight(b, bSpace)
	return
}
