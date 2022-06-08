package traceID

type AssertInterface interface {
	DebugAssert(msg string) RecordInterface
	InfoAssert(msg string) RecordInterface
	WarnAssert(msg string) RecordInterface
	ErrorAssert(msg string) RecordInterface
	FatalAssert(msg string) RecordInterface
	DebugAssertIf(cond bool, msg string) RecordInterface
	InfoAssertIf(cond bool, msg string) RecordInterface
	WarnAssertIf(cond bool, msg string) RecordInterface
	ErrorAssertIf(cond bool, msg string) RecordInterface
	FatalAssertIf(cond bool, msg string) RecordInterface
}

// ctx group

func (c *Ctx) DebugAssert(msg string) RecordInterface {
	return c.Trace(LevelDebug|LevelAssert, msg)
}

func (c *Ctx) InfoAssert(msg string) RecordInterface {
	return c.Trace(LevelInfo|LevelAssert, msg)
}

func (c *Ctx) WarnAssert(msg string) RecordInterface {
	return c.Trace(LevelWarn|LevelAssert, msg)
}

func (c *Ctx) ErrorAssert(msg string) RecordInterface {
	return c.Trace(LevelError|LevelAssert, msg)
}

func (c *Ctx) FatalAssert(msg string) RecordInterface {
	return c.Trace(LevelFatal|LevelAssert, msg)
}

func (c *Ctx) DebugAssertIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelDebug|LevelAssert, msg)
}

func (c *Ctx) InfoAssertIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelInfo|LevelAssert, msg)
}

func (c *Ctx) WarnAssertIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelWarn|LevelAssert, msg)
}

func (c *Ctx) ErrorAssertIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelError|LevelAssert, msg)
}

func (c *Ctx) FatalAssertIf(cond bool, msg string) RecordInterface {
	return c.TraceIf(cond, LevelFatal|LevelAssert, msg)
}

// ctx group end

// thread group

func (t Thread) DebugAssert(msg string) RecordInterface {
	return t.Trace(LevelDebug|LevelAssert, msg)
}

func (t Thread) InfoAssert(msg string) RecordInterface {
	return t.Trace(LevelInfo|LevelAssert, msg)
}

func (t Thread) WarnAssert(msg string) RecordInterface {
	return t.Trace(LevelWarn|LevelAssert, msg)
}

func (t Thread) ErrorAssert(msg string) RecordInterface {
	return t.Trace(LevelError|LevelAssert, msg)
}

func (t Thread) FatalAssert(msg string) RecordInterface {
	return t.Trace(LevelFatal|LevelAssert, msg)
}

func (t Thread) DebugAssertIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelDebug|LevelAssert, msg)
}

func (t Thread) InfoAssertIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelInfo|LevelAssert, msg)
}

func (t Thread) WarnAssertIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelWarn|LevelAssert, msg)
}

func (t Thread) ErrorAssertIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelError|LevelAssert, msg)
}

func (t Thread) FatalAssertIf(cond bool, msg string) RecordInterface {
	return t.TraceIf(cond, LevelFatal|LevelAssert, msg)
}

// thread group end

type DummyAssert struct{}

func (t DummyAssert) DebugAssert(string) RecordInterface         { return DummyRecord{} }
func (t DummyAssert) InfoAssert(string) RecordInterface          { return DummyRecord{} }
func (t DummyAssert) WarnAssert(string) RecordInterface          { return DummyRecord{} }
func (t DummyAssert) ErrorAssert(string) RecordInterface         { return DummyRecord{} }
func (t DummyAssert) FatalAssert(string) RecordInterface         { return DummyRecord{} }
func (t DummyAssert) DebugAssertIf(bool, string) RecordInterface { return DummyRecord{} }
func (t DummyAssert) InfoAssertIf(bool, string) RecordInterface  { return DummyRecord{} }
func (t DummyAssert) WarnAssertIf(bool, string) RecordInterface  { return DummyRecord{} }
func (t DummyAssert) ErrorAssertIf(bool, string) RecordInterface { return DummyRecord{} }
func (t DummyAssert) FatalAssertIf(bool, string) RecordInterface { return DummyRecord{} }
