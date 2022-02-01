package traceID

type Broadcaster interface {
	Broadcast([]byte) (int, error)
}

var (
	bcs []Broadcaster
)

func RegisterBroadcaster(bc Broadcaster) {
	bcs = append(bcs, bc)
}
