package listener

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HTTP struct{}

func (l HTTP) Listen(addr string, out chan []byte) error {
	var (
		u   *url.URL
		err error
	)
	u, err = url.Parse(addr)
	if err != nil {
		return err
	}
	var a string
	if len(u.Scheme) > 0 {
		a += u.Scheme + "://"
	}
	a += u.Host
	p := u.Path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	http.HandleFunc(p, func(w http.ResponseWriter, req *http.Request) {
		p, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("err '%s' caught on request '%s'\n", err.Error(), req.RequestURI)
		} else {
			out <- p
		}
	})
	return http.ListenAndServe(a, nil)
}