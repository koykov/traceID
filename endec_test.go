package traceID

import (
	"bytes"
	"testing"
)

func TestEndec(t *testing.T) {
	var encoded = []byte{
		0x01, 0x00, 0x10, 0x00, 0x48, 0x38, 0x62, 0x71, 0x63, 0x34, 0x71, 0x47, 0x57, 0x71, 0x65, 0x34,
		0x32, 0x6d, 0x62, 0x33, 0x04, 0x00, 0x00, 0xc8, 0x0f, 0x5f, 0xa3, 0x1c, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00,
		0xc8, 0x0f, 0x5f, 0xa3, 0x1c, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x00, 0x00,
		0x1b, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0xc8, 0x0f, 0x5f, 0xa3, 0x1c, 0x00, 0x00,
		0x00, 0x24, 0x00, 0x00, 0x00, 0x1b, 0x00, 0x00, 0x00, 0x2d, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00,
		0x00, 0x00, 0xc8, 0x0f, 0x5f, 0xa3, 0x1c, 0x00, 0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x2d, 0x00,
		0x00, 0x00, 0x71, 0x00, 0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x71, 0x00, 0x00, 0x00, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x5f, 0x31, 0x32, 0x0a, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
		0x5f, 0x32, 0x33, 0x2e, 0x31, 0x34, 0x31, 0x35, 0x0a, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
		0x5f, 0x33, 0x22, 0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72, 0x22, 0x0a, 0x65, 0x78, 0x61, 0x6d, 0x70,
		0x6c, 0x65, 0x5f, 0x34, 0x7b, 0x22, 0x61, 0x22, 0x3a, 0x31, 0x32, 0x33, 0x2c, 0x22, 0x62, 0x22,
		0x3a, 0x34, 0x35, 0x36, 0x2e, 0x37, 0x38, 0x39, 0x2c, 0x22, 0x63, 0x22, 0x3a, 0x22, 0x63, 0x58,
		0x64, 0x6c, 0x63, 0x6e, 0x52, 0x35, 0x22, 0x2c, 0x22, 0x64, 0x22, 0x3a, 0x22, 0x61, 0x73, 0x64,
		0x66, 0x67, 0x68, 0x22, 0x2c, 0x22, 0x65, 0x22, 0x3a, 0x74, 0x72, 0x75, 0x65, 0x7d, 0x0a,
	}
	t.Run("encode", func(t *testing.T) {
		ctx := NewCtx()
		ctx.SetClock(DummyClock{}).
			SetID("H8bqc4qGWqe42mb3").
			Log("example_1", 2).
			Log("example_2", 3.1415).
			Log("example_3", "foobar").
			Log("example_4", struct {
				A int32   `json:"a"`
				B float64 `json:"b"`
				C []byte  `json:"c"`
				D string  `json:"d"`
				E bool    `json:"e"`
			}{
				A: 123,
				B: 456.789,
				C: []byte("qwerty"),
				D: "asdfgh",
				E: true,
			})
		cb := Encode(ctx)
		if !bytes.Equal(cb, encoded) {
			t.FailNow()
		}
	})
	t.Run("decode", func(t *testing.T) {
		var x Packet
		if err := Decode(encoded, &x); err != nil {
			t.Error(err)
		}
		if x.ID != "H8bqc4qGWqe42mb3" {
			t.Error("ID mismatch")
		}
		if len(x.Records) != 4 {
			t.FailNow()
		}
		if x.Records[2].Type != EntryLog {
			t.FailNow()
		}
		if x.Records[2].Time != 123000000456 {
			t.FailNow()
		}
		if lo, hi := x.Records[2].Key.Decode(); string(x.Buf[lo:hi]) != "example_3" {
			t.FailNow()
		}
		if lo, hi := x.Records[2].Value.Decode(); string(x.Buf[lo:hi]) != "\"foobar\"\n" {
			t.FailNow()
		}
	})
}
