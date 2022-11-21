package gin

const (
	request direction = iota
	response
)

type direction int8

func (r direction) String() string {
	switch r {
	case request:
		return "-->"
	case response:
		return "<--"
	default:
		panic("unknown relation type")
	}
}
