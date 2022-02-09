package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB string `json:"db"`
	UI string `json:"ui"`

	Listeners []Listener `json:"listeners"`
	Notifiers []Notifier `json:"notifiers"`
}

type Listener struct {
	Handler string `json:"handler"`
	Addr    string `json:"addr"`
}

type Notifier struct {
	Handler string `json:"handler"`
	Addr    string `json:"addr"`
}

func ParseConfig(filepath string) (*Config, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var c Config
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
