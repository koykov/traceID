package marshaller

import (
	"encoding/json"
	"io"

	"github.com/koykov/bytealg"
)

type JSON struct{}

func (m JSON) Marshal(rw io.ReadWriter, x interface{}, indent bool) (b []byte, err error) {
	e := json.NewEncoder(rw)
	if indent {
		e.SetIndent("", "\t")
	}
	if err = e.Encode(x); err != nil {
		return
	}
	b, err = io.ReadAll(rw)
	b = bytealg.TrimRight(b, bNl)
	return
}
