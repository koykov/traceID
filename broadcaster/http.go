package broadcaster

import (
	"bytes"
	"context"
	"net/http"
	"sync"
)

type HTTP struct {
	Addr string
	mux  sync.Mutex
	buf  bytes.Buffer
}

func (b *HTTP) Broadcast(ctx context.Context, p []byte) (n int, err error) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.buf.Reset()
	if _, err = b.buf.Write(p); err != nil {
		return
	}
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, "GET", b.Addr, &b.buf); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	_, err = http.DefaultClient.Do(req)
	n = len(p)
	return
}
