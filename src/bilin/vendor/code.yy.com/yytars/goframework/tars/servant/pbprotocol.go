// @author kordenlu
// @创建时间 2018/02/05 14:59
// 功能描述: proto buffer support

package servant

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"math"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"context"
)

type PbDispatcher interface {
	Dispatch(context.Context,interface{}, *pbtaf.RequestPacket) (*pbtaf.ResponsePacket, error)
}

type PbProtocol struct {
	dispatcher PbDispatcher
	serverImp  interface{}
}

func NewPbProtocol(dispatcher PbDispatcher, imp interface{}) *PbProtocol {
	s := &PbProtocol{dispatcher: dispatcher, serverImp: imp}
	return s
}

func (s *PbProtocol) Invoke(ctx context.Context,req []byte) ([]byte, error) {
	defer checkPanic()
	var (
		reqPackage pbtaf.RequestPacket
		rspPackage *pbtaf.ResponsePacket
		err        error
		length     uint
	)

	if err := proto.Unmarshal(req, &reqPackage); err != nil {
		appzaplog.Error("Unmarshal req failed", zap.Error(err))
		return nil, err
	}

	rspPackage, err = s.dispatcher.Dispatch(ctx,s.serverImp, &reqPackage)
	if err != nil || rspPackage == nil {
		return nil, err
	}

	rsp, err := proto.Marshal(rspPackage)
	if err != nil {
		appzaplog.Error("Marshal rspPackage failed", zap.Error(err))
		return nil, err
	}

	const (
		sizeLen = 4
	)

	if length = uint(len(rsp)) + sizeLen; length > math.MaxUint32 {
		return nil, errors.New(fmt.Sprintf("grpc: message too large (%d bytes)", length))
	}

	var buf = make([]byte, length)
	binary.BigEndian.PutUint32(buf[0:], uint32(length))
	copy(buf[4:], rsp)
	return buf, nil
}

func (s *PbProtocol) ParsePackage(buff []byte) (int, int) {
	return TafRequest(buff)
}

func (s *PbProtocol) InvokeTimeout(pkg []byte) ([]byte, error) {
	//TODO ERROR PACKAGE
	appzaplog.Error("invoke timeout", zap.Binary("pkg", pkg))
	payload := []byte("timeout")
	ret := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(pkg[:4], uint32(len(ret)))
	copy(pkg[4:], payload)
	return ret, nil
}
