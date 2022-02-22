package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/koykov/bytealg"
)

type Slack struct {
	notifier
	channel, username, template, url string
}

func (n Slack) Notify(ctx context.Context, id string) (err error) {
	if n.channel, n.username, n.template, n.url, err = n.parseAddr(); err != nil {
		return
	}

	msg := strings.ReplaceAll(n.template, "{TID}", id)

	x := struct {
		Channel  string `json:"channel"`
		Username string `json:"username"`
		Text     string `json:"text"`
	}{
		Channel:  n.channel,
		Username: n.username,
		Text:     msg,
	}
	payload, _ := json.Marshal(x)

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, "POST", n.url, bytes.NewBuffer(payload)); err != nil {
		return
	}
	_, err = http.DefaultClient.Do(req)

	return
}

func (n Slack) parseAddr() (channel, username, template, url string, err error) {
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
		case "template":
			template = v
		}
	}
	return
}
