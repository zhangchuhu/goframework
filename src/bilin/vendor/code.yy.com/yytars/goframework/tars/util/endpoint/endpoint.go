package endpoint

type Endpoint struct {
	Host      string
	Port      int32
	IPPort    string // Host:Port
	Timeout   int32
	Proto     string
	Bind      string
	Container string
}

func (ed *Endpoint) istcp() int32 {
	var ret int32
	if ed.Proto == "tcp" {
		ret = 1
	}
	return ret
}