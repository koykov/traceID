package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Telegram struct {
	notifier
}

func (n Telegram) Notify(ctx context.Context, id string) (err error) {
	x := struct {
		ChatID    string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode,omitempty"`
	}{
		ChatID: n.conf.ChatID,
		Text:   n.render(id),
	}
	if len(n.conf.Format) > 0 {
		x.ParseMode = n.conf.Format
	}
	payload, _ := json.Marshal(x)

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, "POST", n.conf.Addr, bytes.NewBuffer(payload)); err != nil {
		return
	}
	req.Header.Add("Host", "api.telegram.org")
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)

	return
}
