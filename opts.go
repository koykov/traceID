package traceID

type Options struct {
	Marshaller Marshaller
	Indent     bool
}

var (
	_, _ = WithMarshaller, WithIndent
)

func WithMarshaller(m8r Marshaller) *Options {
	o := Options{Marshaller: m8r}
	return &o
}

func WithIndent(ind bool) *Options {
	o := Options{Indent: ind}
	return &o
}

func (o *Options) WithMarshaller(m8r Marshaller) *Options {
	o.Marshaller = m8r
	return o
}

func (o *Options) WithIndent(ind bool) *Options {
	o.Indent = ind
	return o
}
