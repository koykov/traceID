package marshaller

import (
	"encoding/json"
	"io"

	"github.com/koykov/bytealg"
)

type JSONFmt struct{}

func (m JSONFmt) Marshal(rw io.ReadWriter, x interface{}) (b []byte, err error) {
	e := json.NewEncoder(rw)
	e.SetIndent("", "\t")
	if err = e.Encode(x); err != nil {
		return
	}
	b, err = io.ReadAll(rw)
	b = bytealg.TrimRight(b, []byte{'\n'})
	return
}
