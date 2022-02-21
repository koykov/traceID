package notifier

type notifier struct {
	addr string
}

func (l *notifier) SetAddr(addr string) {
	l.addr = addr
}
