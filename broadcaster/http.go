package broadcaster

import (
	"bytes"
	"net/http"
	"sync"
)

type HTTP struct {
	Addr string
	mux  sync.Mutex
	buf  bytes.Buffer
}

func (b *HTTP) Broadcast(p []byte) (n int, err error) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.buf.Reset()
	if _, err = b.buf.Write(p); err != nil {
		return
	}
	if _, err = http.Post(b.Addr, "application/octet-stream", &b.buf); err != nil {
		return
	}
	n = len(p)
	return
}
