package traceID

const (
	DefaultZeroMQTopic  = "TRACE_NATIVE"
	ProtobufZeroMQTopic = "TRACE_PROTOBUF"
	DefaultZeroMQHWM    = 1000
)

var (
	// todo implement tracepb and remove me
	_ = ProtobufZeroMQTopic
)
