package listener

import (
	"context"
	"io"
	"log"
	"net/http"
)

type HTTP struct {
	Base
	srv *http.Server
}

func (l HTTP) Listen(ctx context.Context, out chan []byte) (err error) {
	l.srv = &http.Server{Addr: l.conf.Addr}

	http.HandleFunc(l.conf.Path, func(w http.ResponseWriter, req *http.Request) {
		p, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("err '%s' caught on request '%s'\n", err.Error(), req.RequestURI)
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			out <- p
		}
		w.WriteHeader(http.StatusAccepted)
	})

	go func() {
		if err := l.srv.ListenAndServe(); err != nil {
			log.Printf("server '%s' failed with err '%s'\n", l.conf.Addr, err.Error())
		}
	}()

	select {
	case <-ctx.Done():
		err = l.srv.Shutdown(context.Background())
	}
	return
}
