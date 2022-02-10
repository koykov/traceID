package listener

type listener struct {
	addr string
}

func (l *listener) SetAddr(addr string) {
	l.addr = addr
}
