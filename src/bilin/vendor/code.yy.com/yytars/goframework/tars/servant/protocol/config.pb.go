// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

/*
Package pbtaf is a generated protocol buffer package.

It is generated from these files:
	config.proto
	header.proto

It has these top-level messages:
	ConfigPushNotice
	RequestPacket
	ResponsePacket
*/
package pbtaf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ConfigPushNotice_RETCODE int32

const (
	ConfigPushNotice_SUCCESS ConfigPushNotice_RETCODE = 0
	ConfigPushNotice_FAILED  ConfigPushNotice_RETCODE = 1
)

var ConfigPushNotice_RETCODE_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILED",
}
var ConfigPushNotice_RETCODE_value = map[string]int32{
	"SUCCESS": 0,
	"FAILED":  1,
}

func (x ConfigPushNotice_RETCODE) String() string {
	return proto.EnumName(ConfigPushNotice_RETCODE_name, int32(x))
}
func (ConfigPushNotice_RETCODE) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type ConfigPushNotice struct {
	Filename string                   `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
	Retcode  ConfigPushNotice_RETCODE `protobuf:"varint,2,opt,name=retcode,enum=pbtaf.ConfigPushNotice_RETCODE" json:"retcode,omitempty"`
}

func (m *ConfigPushNotice) Reset()                    { *m = ConfigPushNotice{} }
func (m *ConfigPushNotice) String() string            { return proto.CompactTextString(m) }
func (*ConfigPushNotice) ProtoMessage()               {}
func (*ConfigPushNotice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ConfigPushNotice) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *ConfigPushNotice) GetRetcode() ConfigPushNotice_RETCODE {
	if m != nil {
		return m.Retcode
	}
	return ConfigPushNotice_SUCCESS
}

func init() {
	proto.RegisterType((*ConfigPushNotice)(nil), "pbtaf.ConfigPushNotice")
	proto.RegisterEnum("pbtaf.ConfigPushNotice_RETCODE", ConfigPushNotice_RETCODE_name, ConfigPushNotice_RETCODE_value)
}

func init() { proto.RegisterFile("config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 156 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0x48, 0x2a, 0x49, 0x4c, 0x53,
	0xea, 0x65, 0xe4, 0x12, 0x70, 0x06, 0x8b, 0x07, 0x94, 0x16, 0x67, 0xf8, 0xe5, 0x97, 0x64, 0x26,
	0xa7, 0x0a, 0x49, 0x71, 0x71, 0xa4, 0x65, 0xe6, 0xa4, 0xe6, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x06, 0xc1, 0xf9, 0x42, 0x96, 0x5c, 0xec, 0x45, 0xa9, 0x25, 0xc9, 0xf9, 0x29,
	0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x7c, 0x46, 0xf2, 0x7a, 0x60, 0x93, 0xf4, 0xd0, 0x4d, 0xd1,
	0x0b, 0x72, 0x0d, 0x71, 0xf6, 0x77, 0x71, 0x0d, 0x82, 0xa9, 0x57, 0x52, 0xe2, 0x62, 0x87, 0x8a,
	0x09, 0x71, 0x73, 0xb1, 0x07, 0x87, 0x3a, 0x3b, 0xbb, 0x06, 0x07, 0x0b, 0x30, 0x08, 0x71, 0x71,
	0xb1, 0xb9, 0x39, 0x7a, 0xfa, 0xb8, 0xba, 0x08, 0x30, 0x26, 0xb1, 0x81, 0x5d, 0x67, 0x0c, 0x08,
	0x00, 0x00, 0xff, 0xff, 0xb6, 0xc7, 0xcd, 0xb4, 0xad, 0x00, 0x00, 0x00,
}