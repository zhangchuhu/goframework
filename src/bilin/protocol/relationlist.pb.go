// Code generated by protoc-gen-tars. DO NOT EDIT.
// source: relationlist.proto

package bilin

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	"code.yy.com/yytars/goframework/tars/servant"
	"code.yy.com/yytars/goframework/tars/servant/model"
	"code.yy.com/yytars/goframework/jce/taf"
	"errors"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RSUserMikeOptionReq_MIKEOPT int32

const (
	RSUserMikeOptionReq_UNMIKE RSUserMikeOptionReq_MIKEOPT = 0
	RSUserMikeOptionReq_ONMIKE RSUserMikeOptionReq_MIKEOPT = 1
)

var RSUserMikeOptionReq_MIKEOPT_name = map[int32]string{
	0: "UNMIKE",
	1: "ONMIKE",
}
var RSUserMikeOptionReq_MIKEOPT_value = map[string]int32{
	"UNMIKE": 0,
	"ONMIKE": 1,
}

func (x RSUserMikeOptionReq_MIKEOPT) String() string {
	return proto.EnumName(RSUserMikeOptionReq_MIKEOPT_name, int32(x))
}
func (RSUserMikeOptionReq_MIKEOPT) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor14, []int{1, 0}
}

type RelationInfo struct {
	Bilinid       uint64 `protobuf:"varint,1,opt,name=bilinid" json:"bilinid,omitempty"`
	Nick          string `protobuf:"bytes,2,opt,name=nick" json:"nick,omitempty"`
	Avatar        string `protobuf:"bytes,3,opt,name=avatar" json:"avatar,omitempty"`
	Headgear      string `protobuf:"bytes,4,opt,name=headgear" json:"headgear,omitempty"`
	Relationvalue uint64 `protobuf:"varint,5,opt,name=relationvalue" json:"relationvalue,omitempty"`
}

func (m *RelationInfo) Reset()                    { *m = RelationInfo{} }
func (m *RelationInfo) String() string            { return proto.CompactTextString(m) }
func (*RelationInfo) ProtoMessage()               {}
func (*RelationInfo) Descriptor() ([]byte, []int) { return fileDescriptor14, []int{0} }

func (m *RelationInfo) GetBilinid() uint64 {
	if m != nil {
		return m.Bilinid
	}
	return 0
}

func (m *RelationInfo) GetNick() string {
	if m != nil {
		return m.Nick
	}
	return ""
}

func (m *RelationInfo) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *RelationInfo) GetHeadgear() string {
	if m != nil {
		return m.Headgear
	}
	return ""
}

func (m *RelationInfo) GetRelationvalue() uint64 {
	if m != nil {
		return m.Relationvalue
	}
	return 0
}

// 用户上下麦操作
type RSUserMikeOptionReq struct {
	Header *Header                     `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Owner  uint64                      `protobuf:"varint,2,opt,name=owner" json:"owner,omitempty"`
	Opt    RSUserMikeOptionReq_MIKEOPT `protobuf:"varint,3,opt,name=opt,enum=bilin.relationlist.RSUserMikeOptionReq_MIKEOPT" json:"opt,omitempty"`
}

func (m *RSUserMikeOptionReq) Reset()                    { *m = RSUserMikeOptionReq{} }
func (m *RSUserMikeOptionReq) String() string            { return proto.CompactTextString(m) }
func (*RSUserMikeOptionReq) ProtoMessage()               {}
func (*RSUserMikeOptionReq) Descriptor() ([]byte, []int) { return fileDescriptor14, []int{1} }

func (m *RSUserMikeOptionReq) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *RSUserMikeOptionReq) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

func (m *RSUserMikeOptionReq) GetOpt() RSUserMikeOptionReq_MIKEOPT {
	if m != nil {
		return m.Opt
	}
	return RSUserMikeOptionReq_UNMIKE
}

type RSUserMikeOptionResp struct {
	Commonret *CommonRetInfo `protobuf:"bytes,1,opt,name=commonret" json:"commonret,omitempty"`
}

func (m *RSUserMikeOptionResp) Reset()                    { *m = RSUserMikeOptionResp{} }
func (m *RSUserMikeOptionResp) String() string            { return proto.CompactTextString(m) }
func (*RSUserMikeOptionResp) ProtoMessage()               {}
func (*RSUserMikeOptionResp) Descriptor() ([]byte, []int) { return fileDescriptor14, []int{2} }

func (m *RSUserMikeOptionResp) GetCommonret() *CommonRetInfo {
	if m != nil {
		return m.Commonret
	}
	return nil
}

type GetUserRelationMedalReq struct {
	Header *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Owner  uint64  `protobuf:"varint,2,opt,name=owner" json:"owner,omitempty"`
}

func (m *GetUserRelationMedalReq) Reset()                    { *m = GetUserRelationMedalReq{} }
func (m *GetUserRelationMedalReq) String() string            { return proto.CompactTextString(m) }
func (*GetUserRelationMedalReq) ProtoMessage()               {}
func (*GetUserRelationMedalReq) Descriptor() ([]byte, []int) { return fileDescriptor14, []int{3} }

func (m *GetUserRelationMedalReq) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetUserRelationMedalReq) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

type GetUserRelationMedalResp struct {
	Commonret *CommonRetInfo `protobuf:"bytes,1,opt,name=commonret" json:"commonret,omitempty"`
	Medalid   uint32         `protobuf:"varint,2,opt,name=medalid" json:"medalid,omitempty"`
	Medalname string         `protobuf:"bytes,3,opt,name=medalname" json:"medalname,omitempty"`
	MedalUrl  string         `protobuf:"bytes,4,opt,name=medalUrl" json:"medalUrl,omitempty"`
}

func (m *GetUserRelationMedalResp) Reset()                    { *m = GetUserRelationMedalResp{} }
func (m *GetUserRelationMedalResp) String() string            { return proto.CompactTextString(m) }
func (*GetUserRelationMedalResp) ProtoMessage()               {}
func (*GetUserRelationMedalResp) Descriptor() ([]byte, []int) { return fileDescriptor14, []int{4} }

func (m *GetUserRelationMedalResp) GetCommonret() *CommonRetInfo {
	if m != nil {
		return m.Commonret
	}
	return nil
}

func (m *GetUserRelationMedalResp) GetMedalid() uint32 {
	if m != nil {
		return m.Medalid
	}
	return 0
}

func (m *GetUserRelationMedalResp) GetMedalname() string {
	if m != nil {
		return m.Medalname
	}
	return ""
}

func (m *GetUserRelationMedalResp) GetMedalUrl() string {
	if m != nil {
		return m.MedalUrl
	}
	return ""
}

func init() {
	proto.RegisterType((*RelationInfo)(nil), "bilin.relationlist.RelationInfo")
	proto.RegisterType((*RSUserMikeOptionReq)(nil), "bilin.relationlist.RSUserMikeOptionReq")
	proto.RegisterType((*RSUserMikeOptionResp)(nil), "bilin.relationlist.RSUserMikeOptionResp")
	proto.RegisterType((*GetUserRelationMedalReq)(nil), "bilin.relationlist.GetUserRelationMedalReq")
	proto.RegisterType((*GetUserRelationMedalResp)(nil), "bilin.relationlist.GetUserRelationMedalResp")
	proto.RegisterEnum("bilin.relationlist.RSUserMikeOptionReq_MIKEOPT", RSUserMikeOptionReq_MIKEOPT_name, RSUserMikeOptionReq_MIKEOPT_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context

// Client API for RelationListServant service

type RelationListServantClient interface {
	// 用户上下麦 RS --> Relation Servant
	RSUserMikeOption(ctx context.Context, in *RSUserMikeOptionReq, opts ...map[string]string) (*RSUserMikeOptionResp, error)
	// 获取亲密度勋章数据
	GetUserRelationMedal(ctx context.Context, in *GetUserRelationMedalReq, opts ...map[string]string) (*GetUserRelationMedalResp, error)
}

type relationListServantClient struct {
	s model.Servant
}

func NewRelationListServantClient(objname string, comm servant.ICommunicator) RelationListServantClient {
	if comm == nil || objname == "" {
		return nil
	}
	return &relationListServantClient{s: comm.GetServantProxy(objname)}
}

func (c *relationListServantClient) RSUserMikeOption(ctx context.Context, in *RSUserMikeOptionReq, opts ...map[string]string) (*RSUserMikeOptionResp, error) {
	var (
		reply RSUserMikeOptionResp
	)

	pbbuf, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	_resp, err := c.s.Taf_invoke(ctx, 0, "RSUserMikeOption", pbbuf)
	if err != nil {
		return nil, err
	}

	if err = proto.Unmarshal(_resp.SBuffer, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
func (c *relationListServantClient) GetUserRelationMedal(ctx context.Context, in *GetUserRelationMedalReq, opts ...map[string]string) (*GetUserRelationMedalResp, error) {
	var (
		reply GetUserRelationMedalResp
	)

	pbbuf, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	_resp, err := c.s.Taf_invoke(ctx, 0, "GetUserRelationMedal", pbbuf)
	if err != nil {
		return nil, err
	}

	if err = proto.Unmarshal(_resp.SBuffer, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

// Server API for RelationListServant service

type RelationListServantServer interface {
	// 用户上下麦 RS --> Relation Servant
	RSUserMikeOption(context.Context, *RSUserMikeOptionReq) (*RSUserMikeOptionResp, error)
	// 获取亲密度勋章数据
	GetUserRelationMedal(context.Context, *GetUserRelationMedalReq) (*GetUserRelationMedalResp, error)
}

type relationListServantDispatcher struct {
}

func NewRelationListServantDispatcher() servant.Dispatcher {
	return &relationListServantDispatcher{}
}

func (_obj *relationListServantDispatcher) Dispatch(ctx context.Context, _val interface{}, req *taf.RequestPacket) (*taf.ResponsePacket, error) {
	var pbbuf []byte
	_imp := _val.(RelationListServantServer)
	switch req.SFuncName {
	case "RSUserMikeOption":
		var req_ RSUserMikeOptionReq
		if err := proto.Unmarshal(req.SBuffer, &req_); err != nil {
			return nil, err
		}

		_ret, err := _imp.RSUserMikeOption(ctx, &req_)
		if err != nil {
			return nil, err
		}

		if pbbuf, err = proto.Marshal(_ret); err != nil {
			return nil, err
		}

	case "GetUserRelationMedal":
		var req_ GetUserRelationMedalReq
		if err := proto.Unmarshal(req.SBuffer, &req_); err != nil {
			return nil, err
		}

		_ret, err := _imp.GetUserRelationMedal(ctx, &req_)
		if err != nil {
			return nil, err
		}

		if pbbuf, err = proto.Marshal(_ret); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unknow func")
	}
	return &taf.ResponsePacket{
		IVersion:   1,
		IRequestId: req.IRequestId,
		SBuffer:    pbbuf,
		Context:    req.Context,
	}, nil
}

func init() { proto.RegisterFile("relationlist.proto", fileDescriptor14) }

var fileDescriptor14 = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x51, 0x8e, 0xd3, 0x30,
	0x10, 0x25, 0x6c, 0x9a, 0xd2, 0x61, 0x8b, 0x56, 0x43, 0x04, 0x51, 0xc4, 0xc7, 0x12, 0x81, 0xa8,
	0x04, 0x0a, 0x52, 0x38, 0x01, 0x20, 0x04, 0x0b, 0x94, 0x22, 0x2f, 0xe5, 0x83, 0x3f, 0xef, 0x66,
	0x00, 0xab, 0x89, 0x9d, 0x3a, 0xa6, 0xdc, 0x84, 0x03, 0x70, 0x0a, 0x6e, 0xc5, 0x15, 0x90, 0x1d,
	0x87, 0x02, 0x0d, 0x52, 0x05, 0x7f, 0xf3, 0x66, 0x26, 0x6f, 0x9e, 0x67, 0x5e, 0x00, 0x35, 0x55,
	0xdc, 0x08, 0x25, 0x2b, 0xd1, 0x9a, 0xbc, 0xd1, 0xca, 0x28, 0xc4, 0x33, 0x51, 0x09, 0x99, 0xff,
	0x5a, 0x49, 0x0f, 0x3f, 0x12, 0x2f, 0x49, 0x77, 0x1d, 0xd9, 0x97, 0x00, 0x0e, 0x99, 0x2f, 0x9f,
	0xc8, 0xf7, 0x0a, 0x13, 0x18, 0xbb, 0x8f, 0x44, 0x99, 0x04, 0xc7, 0xc1, 0x2c, 0x64, 0x3d, 0x44,
	0x84, 0x50, 0x8a, 0xf3, 0x55, 0x72, 0xf1, 0x38, 0x98, 0x4d, 0x98, 0x8b, 0xf1, 0x1a, 0x44, 0x7c,
	0xc3, 0x0d, 0xd7, 0xc9, 0x81, 0xcb, 0x7a, 0x84, 0x29, 0x5c, 0xb2, 0x63, 0x3e, 0x10, 0xd7, 0x49,
	0xe8, 0x2a, 0x3f, 0x31, 0xde, 0x82, 0x69, 0x2f, 0x68, 0xc3, 0xab, 0x4f, 0x94, 0x8c, 0xdc, 0x9c,
	0xdf, 0x93, 0xd9, 0xb7, 0x00, 0xae, 0xb2, 0xd3, 0x65, 0x4b, 0x7a, 0x2e, 0x56, 0xb4, 0x68, 0x6c,
	0x85, 0xd1, 0x1a, 0x6f, 0x43, 0xd4, 0x3d, 0xc0, 0xc9, 0xbb, 0x5c, 0x4c, 0xf3, 0xee, 0x8d, 0xcf,
	0x5c, 0x92, 0xf9, 0x22, 0xc6, 0x30, 0x52, 0x9f, 0x25, 0x69, 0xa7, 0x36, 0x64, 0x1d, 0xc0, 0x87,
	0x70, 0xa0, 0x1a, 0xe3, 0xb4, 0x5e, 0x29, 0xee, 0xe7, 0xbb, 0xdb, 0xc9, 0x07, 0x46, 0xe6, 0xf3,
	0x93, 0x17, 0x4f, 0x16, 0xaf, 0xdf, 0x30, 0xfb, 0x6d, 0x76, 0x13, 0xc6, 0x1e, 0x23, 0x40, 0xb4,
	0x7c, 0x65, 0xc1, 0xd1, 0x05, 0x1b, 0x2f, 0xba, 0x38, 0xc8, 0x9e, 0x43, 0xbc, 0x4b, 0xd3, 0x36,
	0x58, 0xc0, 0xe4, 0x5c, 0xd5, 0xb5, 0x92, 0x9a, 0x8c, 0x57, 0x1f, 0x7b, 0x0d, 0x8f, 0x5d, 0x9e,
	0x91, 0xb1, 0x37, 0x60, 0xdb, 0xb6, 0xec, 0x2d, 0x5c, 0x7f, 0x4a, 0xc6, 0x92, 0xf5, 0x57, 0x9a,
	0x53, 0xc9, 0xab, 0xff, 0xdd, 0x44, 0xf6, 0x35, 0x80, 0x64, 0x98, 0xf8, 0xdf, 0x84, 0x5a, 0xdf,
	0xd4, 0x96, 0x40, 0x94, 0x6e, 0xd0, 0x94, 0xf5, 0x10, 0x6f, 0xc0, 0xc4, 0x85, 0x92, 0xd7, 0xe4,
	0x6d, 0xb2, 0x4d, 0x58, 0xa7, 0x38, 0xb0, 0xd4, 0x55, 0xef, 0x94, 0x1e, 0x17, 0xdf, 0xad, 0x07,
	0xbc, 0xba, 0x97, 0xa2, 0x35, 0xa7, 0xa4, 0x37, 0x5c, 0x1a, 0x24, 0x38, 0xfa, 0x73, 0xc1, 0x78,
	0x67, 0xcf, 0x6b, 0xa6, 0xb3, 0xfd, 0x1a, 0xdb, 0x06, 0xd7, 0x10, 0x0f, 0xad, 0x08, 0xef, 0x0e,
	0x31, 0xfc, 0xe5, 0x4a, 0xe9, 0xbd, 0xfd, 0x9b, 0xdb, 0xe6, 0xd1, 0xf8, 0xdd, 0xc8, 0xb5, 0x9f,
	0x45, 0xee, 0xf7, 0x7c, 0xf0, 0x23, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xd4, 0xc8, 0xec, 0xd6, 0x03,
	0x00, 0x00,
}
