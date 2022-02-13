package traceID

import "errors"

var (
	ErrPacketTooShort = errors.New("packet too short")
	ErrHomelessThread = errors.New("homeless thread, use context.AcquireThread() or thread.AcquireThread() to make")
)
