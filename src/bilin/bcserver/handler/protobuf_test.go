package handler

import (
	"bilin/protocol"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestMarshalEmpty(t *testing.T) {
	var err error
	var req bilin.ExitBroRoomReq
	var pbbuf []byte
	req = bilin.ExitBroRoomReq{}
	if pbbuf, err = proto.Marshal(&req); err != nil {
		t.Error(err)
		return
	}
	t.Log("pbbuf value=%v, len=%d", pbbuf, len(pbbuf))
	pbbuf = nil
	t.Log("pbbuf(if nil) value=%v, len=%d", pbbuf, len(pbbuf))
	//t.Error("就是不通过")
}

func TestUnmarshalEmpty(t *testing.T) {
	var err error
	var req bilin.ExitBroRoomReq
	if err = proto.Unmarshal(nil, &req); err != nil {
		t.Error(err)
		return
	}
}
