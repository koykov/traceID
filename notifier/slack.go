package notifier

import (
	"context"
	"strings"

	"github.com/koykov/bytealg"
)

type Slack struct {
	notifier
	channel, username, url string
}

func (n Slack) Notify(ctx context.Context, id string) (err error) {
	_, _ = ctx, id
	return
}

func (n Slack) parseAddr() (channel, username, url string, err error) {
	url = n.addr
	var p int
	for {
		p = strings.Index(url, ";")
		if p == -1 {
			break
		}
		kv := url[:p]
		url = bytealg.TrimStr(url[p+1:], " ")
		p1 := strings.Index(kv, "=")
		if p1 == -1 {
			continue
		}
		k, v := kv[:p1], kv[p1+1:]
		switch k {
		case "channel":
			channel = v
		case "username":
			username = v
		}
	}
	return
}
