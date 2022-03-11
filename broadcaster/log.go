package broadcaster

import (
	"context"
	"io"
	"log"
	"os"
	"sync"

	"github.com/koykov/fastconv"
)

type Log struct {
	Out    io.Writer
	Prefix string
	Flags  int
	once   sync.Once
	inst   *log.Logger
}

func (b *Log) init() {
	if b.Out == nil {
		b.Out = os.Stdout
	}
	if b.Flags == 0 {
		b.Flags = log.LstdFlags
	}
	b.inst = log.New(b.Out, b.Prefix, b.Flags)
}

func (b *Log) Broadcast(ctx context.Context, p []byte) (n int, err error) {
	_ = ctx
	b.once.Do(b.init)
	b.inst.Println(fastconv.B2S(p))
	return
}
