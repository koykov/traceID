package traceID

import "context"

type Notifier interface {
	SetAddr(string)
	Notify(context.Context, string) error
}

type NotifierConfig struct {
	Handler   string `json:"handler"`
	Addr      string `json:"addr,omitempty"`
	Channel   string `json:"channel,omitempty"`
	ChannelID string `json:"channelID,omitempty"`
	ChatID    string `json:"chatID,omitempty"`
	Username  string `json:"username,omitempty"`
	Token     string `json:"token,omitempty"`
	Template  string `json:"template"`
}
