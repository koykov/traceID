package traceID

const (
	OptionMarshaller Option = "marshaller"
	OptionIndent     Option = "indent"
)

type Option string

type optionKV struct {
	k Option
	v interface{}
}
