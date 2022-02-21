package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

type Slack struct {
	notifier
	channel, username, url string
}

func (n Slack) Notify(ctx context.Context, id string) (err error) {
	if n.channel, n.username, n.url, err = n.parseAddr(); err != nil {
		return
	}

	x := struct {
		Channel  string `json:"channel"`
		Username string `json:"username"`
		Text     string `json:"text"`
	}{
		Channel:  n.channel,
		Username: n.username,
		Text:     fmt.Sprintf("New traceID <https://trace.com/%s|#%s> caught.", id, id),
	}
	payload, _ := json.Marshal(x)

	cmd := exec.CommandContext(ctx, "curl", "-X", "POST", "--data-urlencode", fastconv.B2S(payload), n.url)
	err = cmd.Run()

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
