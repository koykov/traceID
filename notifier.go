package traceID

import "context"

type Notifier interface {
	SetConfig(*NotifierConfig)
	Notify(context.Context, string) error
}

type NotifierConfig struct {
	Handler  string `json:"handler"`
	Addr     string `json:"addr,omitempty"`
	Channel  string `json:"channel,omitempty"`
	ChatID   string `json:"chatID,omitempty"`
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
	Template string `json:"template"`
}
