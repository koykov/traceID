package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/koykov/traceID"
)

type worker struct {
	id      uint
	ctx     context.Context
	cancel  context.CancelFunc
	verbose bool
}

type workerRepo struct {
	buf []worker
}

var (
	wsRepo workerRepo
)

func (r *workerRepo) makeWorker(ctx context.Context, cancel context.CancelFunc, verbose bool) *worker {
	w := worker{
		id:      uint(len(r.buf)),
		ctx:     ctx,
		cancel:  cancel,
		verbose: verbose,
	}
	r.buf = append(r.buf, w)
	return &w
}

func (r *workerRepo) stopAll() {
	for i := 0; i < len(r.buf); i++ {
		r.buf[i].cancel()
	}
}

func (w worker) work(bus chan []byte) {
	if w.verbose {
		log.Printf("worker #%d started\n", w.id)
	}
	for {
		select {
		case p := <-bus:
			var msg traceID.Message
			if err := traceID.Decode(p, &msg); err != nil {
				log.Printf("message decode failed: %s\n", err.Error())
				continue
			}
			if w.verbose {
				b, _ := json.Marshal(msg)
				log.Printf("message received: %s", string(b))
			}
			// todo write message to db and notify
		case <-w.ctx.Done():
			if w.verbose {
				log.Printf("worker #%d stopped\n", w.id)
			}
			return
		}
	}
}
