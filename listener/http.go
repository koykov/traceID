package listener

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type HTTP struct {
	listener
	srv *http.Server
}

func (l HTTP) Listen(ctx context.Context, out chan []byte) error {
	host, path, err := l.parseAddr()
	if err != nil {
		return err
	}

	l.srv = &http.Server{Addr: host}

	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
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
			log.Printf("server '%s' failed with err '%s'\n", host, err.Error())
		}
	}()

	select {
	case <-ctx.Done():
		err = l.srv.Shutdown(context.Background())
	}
	return err
}

func (l HTTP) parseAddr() (host, path string, err error) {
	if len(l.addr) == 0 {
		err = errors.New("empty address")
		return
	}
	sli := strings.IndexByte(l.addr, '/')
	if sli == -1 {
		err = errors.New("bad address format")
		return
	}
	host, path = l.addr[:sli], l.addr[sli:]
	return
}
