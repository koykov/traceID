package notifier

import (
	"strings"

	"github.com/koykov/traceID"
)

const (
	defaultTemplate = "New traceID #{TID}"
)

type notifier struct {
	conf *traceID.NotifierConfig
}

func (n *notifier) SetConfig(conf *traceID.NotifierConfig) {
	n.conf = conf
}

func (n *notifier) render(id string) string {
	tpl := defaultTemplate
	if len(n.conf.Template) > 0 {
		tpl = n.conf.Template
	}
	return strings.ReplaceAll(tpl, "{TID}", id)
}
