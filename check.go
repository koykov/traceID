package traceID

import "fmt"

type BytesContainer interface {
	Bytes() []byte
}

type ProtoContainer interface {
	fmt.GoStringer
	Size() int
	Marshal() ([]byte, error)
	MarshalTo([]byte) (int, error)
	MarshalToSizedBuffer([]byte) (int, error)
	Unmarshal([]byte) error
}

func checkStringer(x interface{}) bool {
	_, ok := x.(fmt.Stringer)
	return ok && !checkProto(x)
}

func checkProto(x interface{}) bool {
	_, ok := x.(ProtoContainer)
	return ok
}

func checkBytes(x interface{}) bool {
	_, ok := x.(BytesContainer)
	return ok
}
