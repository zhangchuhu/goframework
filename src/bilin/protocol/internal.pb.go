// Code generated by protoc-gen-tars. DO NOT EDIT.
// source: internal.proto

package bilin

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MultiPush struct {
	Msg     *ServerPush `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	UserIDs []int64     `protobuf:"varint,2,rep,packed,name=userIDs" json:"userIDs,omitempty"`
	AppID   int32       `protobuf:"varint,3,opt,name=appID" json:"appID,omitempty"`
}

func (m *MultiPush) Reset()                    { *m = MultiPush{} }
func (m *MultiPush) String() string            { return proto.CompactTextString(m) }
func (*MultiPush) ProtoMessage()               {}
func (*MultiPush) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

func (m *MultiPush) GetMsg() *ServerPush {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *MultiPush) GetUserIDs() []int64 {
	if m != nil {
		return m.UserIDs
	}
	return nil
}

func (m *MultiPush) GetAppID() int32 {
	if m != nil {
		return m.AppID
	}
	return 0
}

func init() {
	proto.RegisterType((*MultiPush)(nil), "bilin.MultiPush")
}

func init() { proto.RegisterFile("internal.proto", fileDescriptor10) }

var fileDescriptor10 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xcc, 0x2b, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4d, 0xca, 0xcc, 0xc9,
	0xcc, 0x93, 0xe2, 0x2a, 0x28, 0x2d, 0xce, 0x80, 0x08, 0x29, 0x25, 0x70, 0x71, 0xfa, 0x96, 0xe6,
	0x94, 0x64, 0x06, 0x94, 0x16, 0x67, 0x08, 0x29, 0x73, 0x31, 0xe7, 0x16, 0xa7, 0x4b, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x1b, 0x09, 0xea, 0x81, 0x55, 0xeb, 0x05, 0xa7, 0x16, 0x95, 0xa5, 0x16, 0x81,
	0xe4, 0x83, 0x40, 0xb2, 0x42, 0x12, 0x5c, 0xec, 0xa5, 0xc5, 0xa9, 0x45, 0x9e, 0x2e, 0xc5, 0x12,
	0x4c, 0x0a, 0xcc, 0x1a, 0xcc, 0x41, 0x30, 0xae, 0x90, 0x08, 0x17, 0x6b, 0x62, 0x41, 0x81, 0xa7,
	0x8b, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x6b, 0x10, 0x84, 0x93, 0xc4, 0x06, 0xb6, 0xc8, 0x18, 0x10,
	0x00, 0x00, 0xff, 0xff, 0x7c, 0x07, 0xfb, 0xc2, 0x8d, 0x00, 0x00, 0x00,
}
