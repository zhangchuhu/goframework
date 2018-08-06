package handler

import (
	"bilin/protocol"
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/jsonpb"
)

func TestMarshalSearchRsp(t *testing.T) {
	const jsonStr = `{
  "responseHeader": {
    "status": 0,
    "QTime": 11,
	"correctWord": ""
  },
  "response": {
    "1": {
      "numFound": 100,
      "start": 0,
      "docs": [
        {
          "id": "870739132"
        },
        {
          "id": "870739133"
        }
      ]
    },
	"2": {
      "docs": [
        {
          "id": "abc",
          "name": "我们"
        }
      ]
    }
  }
}`
	var rsp InternalSearchRsp
	if err := json.Unmarshal([]byte(jsonStr), &rsp); err != nil {
		t.Fatalf("json unmarshal fail: %v", err)
	}
	rspStr, _ := json.Marshal(&rsp)
	t.Logf("rsp: %s", rspStr)
}

func TestSearch(t *testing.T) {
	rsp, err := NewSearchServantObj().Search(context.Background(), &bilin.SearchReq{
		Q:     `播`,
		Typ:   bilin.SearchType_ROOM,
		Rows:  10,
		Start: 0,
		Uid:   "40373825",
	})
	if err != nil {
		t.Fatalf("call fail: %v", err)
	}
	m := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}
	rspStr, _ := m.MarshalToString(rsp)
	t.Logf("rsp: %s", rspStr)
}

func TestSearchServantObj_GetHotSongs(t *testing.T) {
	rsp, err := NewSearchServantObj().GetHotSongs(context.Background(), &bilin.GetHotSongsReq{
		Rows:  10,
		Start: 140,
	})
	if err != nil {
		t.Fatalf("call fail: %v", err)
	}
	m := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}
	rspStr, _ := m.MarshalToString(rsp)
	t.Logf("rsp: %s", rspStr)
}

func TestSearchServantObj_GetRelatedHotSearches(t *testing.T) {
	rsp, err := NewSearchServantObj().GetRelatedHotSearches(context.Background(), &bilin.GetRelatedHotSearchesReq{
		Q:   `小`,
		Typ: bilin.SearchType_USER_ROOM,
	})
	if err != nil {
		t.Fatalf("call fail: %v", err)
	}
	m := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}
	rspStr, _ := m.MarshalToString(rsp)
	t.Logf("rsp: %s", rspStr)
}
