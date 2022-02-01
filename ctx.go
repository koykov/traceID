package traceID

type Ctx struct {
	id   string
	buf  []kv
	bufl uint
}

type kv struct {
	key string
	val interface{}
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
}

func (c *Ctx) SetID(id string) {
	c.id = id
}

func (c *Ctx) BeginTXN() Interface {
	c.Reset()
	return c
}

func (c *Ctx) Add(val interface{}) Interface {
	return c.AddKV("", val)
}

func (c *Ctx) AddKV(key string, val interface{}) Interface {
	c.buf = append(c.buf, kv{key: key, val: val})
	return c
}

func (c *Ctx) Commit() error {
	// ...
	return nil
}

func (c *Ctx) CommitTo(bc Broadcaster) error {
	// ...
	c.Reset()
	return nil
}

func (c *Ctx) Reset() *Ctx {
	c.buf = c.buf[:0]
	return c
}
