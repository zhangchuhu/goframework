package gojce

//jce type
const (
	BYTE = iota
	SHORT
	INT
	LONG
	FLOAT
	DOUBLE
	STRING1
	STRING4
	MAP
	LIST
	STRUCT_BEGIN
	STRUCT_END
	ZERO_TAG
	SIMPLE_LIST
)

const (
	JCE_MAX_STRING_LENGTH = 100 * 1024 * 1024
)

type head_data struct {
	ty  byte
	tag int
}

func (h *head_data) clear() {
	h.ty = 0
	h.tag = 0
}
