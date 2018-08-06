package servant

import (
	"time"
	base "code.yy.com/yytars/goframework/jce/servant/taf"
	"code.yy.com/yytars/goframework/jce/taf"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
)

type IMessage interface {
	getAdapterProxy() *AdapterProxy
	getFuncName() string
	getSServantName() string
	getRespRet() int32
	Cost() int64
	hashEnable() bool
	consistHashEnable() bool
	HashCode() string
}

type Message struct {
	Req  *taf.RequestPacket
	Resp *taf.ResponsePacket

	Obj *ObjectProxy
	Ser *ServantProxy
	Adp *AdapterProxy

	BeginTime int64
	EndTime   int64
	Status    int

	hashCode string
	isHash   bool
	isConsistHash bool
}

func (m *Message) Init() {
	m.BeginTime = time.Now().UnixNano() / 1000000
}

func (m *Message) End() {
	m.Status = base.JCESERVERSUCCESS
	m.EndTime = time.Now().UnixNano() / 1000000
}

func (m *Message) Cost() int64 {
	return m.EndTime - m.BeginTime
}

func (m *Message) SetHashCode(code string) {
	m.hashCode = code
	m.isHash = true
}

func (m *Message) setConsistHashCode(code string) {
	m.hashCode = code
	m.isConsistHash = true
}

func (m *Message) getAdapterProxy() *AdapterProxy {
	return m.Adp
}

func (m *Message) getFuncName() string {
	return m.Req.SFuncName
}

func (m *Message) getSServantName() string {
	return m.Req.SServantName
}

func (m *Message) getRespRet() int32 {
	if m.Resp != nil {
		return m.Resp.IRet
	}
	return -1
}

func (m *Message) hashEnable() bool {
	return m.isHash
}

func (m *Message)consistHashEnable()bool  {
	return m.isConsistHash
}

func (m *Message) HashCode() string {
	return m.hashCode
}

// protobuf msg
type PbMessage struct {
	Req  *pbtaf.RequestPacket
	Resp *pbtaf.ResponsePacket

	Obj *ObjectProxy
	Ser *ServantProxy
	Adp *AdapterProxy

	BeginTime int64
	EndTime   int64
	Status    int

	hashCode string
	isHash   bool
	isConsistHash bool
}

func (m *PbMessage) Init() {
	m.BeginTime = time.Now().UnixNano() / 1000000
}

func (m *PbMessage) End() {
	m.Status = base.JCESERVERSUCCESS
	m.EndTime = time.Now().UnixNano() / 1000000
}

func (m *PbMessage) Cost() int64 {
	return m.EndTime - m.BeginTime
}

func (m *PbMessage) SetHashCode(code string) {
	m.hashCode = code
	m.isHash = true
}

func (m *PbMessage) getAdapterProxy() *AdapterProxy {
	return m.Adp
}

func (m *PbMessage) getFuncName() string {
	return m.Req.SFuncName
}

func (m *PbMessage) getSServantName() string {
	return m.Req.SServantName
}
func (m *PbMessage) getRespRet() int32 {
	if m.Resp != nil {
		return m.Resp.IRet
	}
	return -1
}

func (m *PbMessage) hashEnable() bool {
	return m.isHash
}

func (m *PbMessage)consistHashEnable()bool  {
	return m.isConsistHash
}

func (m *PbMessage) HashCode() string {
	return m.hashCode
}
