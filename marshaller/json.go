package marshaller

import (
	"encoding/json"
	"io"

	"github.com/koykov/bytealg"
)

type JSON struct{}

func (m JSON) Marshal(rw io.ReadWriter, x interface{}) (b []byte, err error) {
	e := json.NewEncoder(rw)
	if err = e.Encode(x); err != nil {
		return
	}
	b, err = io.ReadAll(rw)
	b = bytealg.TrimRight(b, []byte{'\n'})
	return
}
