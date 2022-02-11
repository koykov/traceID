package traceID

import "errors"

var (
	ErrPacketTooShort = errors.New("packet too short")
	ErrHomelessThread = errors.New("homeless thread, use context.Thread() to make")
)
