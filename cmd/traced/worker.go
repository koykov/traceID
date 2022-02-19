package main

import (
	"context"
	"log"

	"github.com/koykov/traceID"
)

type worker struct {
	id     uint
	ctx    context.Context
	cancel context.CancelFunc
}

type workerRepo struct {
	buf []worker
}

var (
	wsRepo workerRepo
)

func (r *workerRepo) makeWorker(ctx context.Context, cancel context.CancelFunc) {
	w := worker{
		id:     uint(len(r.buf)),
		ctx:    ctx,
		cancel: cancel,
	}
	r.buf = append(r.buf, w)
}

func (r *workerRepo) startWorker(idx uint, bus chan []byte) {
	if idx < uint(len(r.buf)) {
		r.buf[idx].work(bus)
	}
}

func (r *workerRepo) stopAll() {
	for i := 0; i < len(r.buf); i++ {
		r.buf[i].cancel()
	}
}

func (w worker) work(bus chan []byte) {
	for {
		select {
		case p := <-bus:
			var msg traceID.Message
			if err := traceID.Decode(p, &msg); err != nil {
				log.Printf("message decode failed: %s\n", err.Error())
				continue
			}
			// todo write message to db and notify
		case <-w.ctx.Done():
			return
		}
	}
}
