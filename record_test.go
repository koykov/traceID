package traceID

import (
	"bytes"
	"testing"

	"github.com/koykov/traceID/marshaller"
)

func TestRecord(t *testing.T) {
	t.Run("skipWith", func(t *testing.T) {
		ctx := NewCtx()
		ctx.Debug("foobar").
			Var("foo", "some string").
			VarIf(false, "bar", "qwe").With(OptionMarshaller, marshaller.Binary{})
		b := Encode(ctx)
		if !bytes.Contains(b, []byte("some string")) {
			t.FailNow()
		}
	})
}
