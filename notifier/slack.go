package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Slack struct {
	notifier
}

func (n Slack) Notify(ctx context.Context, id string) (err error) {
	x := struct {
		Channel  string `json:"channel"`
		Username string `json:"username"`
		Text     string `json:"text"`
	}{
		Channel:  n.conf.Channel,
		Username: n.conf.Username,
		Text:     n.render(id),
	}
	payload, _ := json.Marshal(x)

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, "POST", n.conf.Addr, bytes.NewBuffer(payload)); err != nil {
		return
	}
	_, err = http.DefaultClient.Do(req)

	return
}
