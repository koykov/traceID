package notifier

import (
	"context"
	"strings"

	"github.com/koykov/bytealg"
)

type Telegram struct {
	notifier
	token, chatID, template string
}

func (n Telegram) Notify(ctx context.Context, id string) (err error) {
	_, _ = ctx, id
	// ...
	return
}

func (n *Telegram) parseAddr() (err error) {
	addr := n.addr
	var (
		p   int
		brk bool
	)
	for {
		p = strings.Index(addr, ";")
		if p == -1 {
			brk = true
			p = len(addr)
		}
		kv := addr[:p]
		if !brk {
			addr = bytealg.TrimStr(addr[p+1:], " ")
		}
		p1 := strings.Index(kv, "=")
		if p1 == -1 {
			continue
		}
		k, v := kv[:p1], kv[p1+1:]
		switch k {
		case "chatID":
			n.chatID = v
		case "token":
			n.token = v
		case "template":
			n.template = v
		}
		if brk {
			break
		}
	}
	return
}
