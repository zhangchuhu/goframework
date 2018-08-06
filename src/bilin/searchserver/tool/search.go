package main

import (
	"bilin/protocol"
	"context"
	"fmt"
	"log"
	"os"

	"code.yy.com/yytars/goframework/tars/servant"
	"github.com/golang/protobuf/jsonpb"
)

func main() {
	comm := servant.NewPbCommunicator()
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 183.36.111.61 -p 17890")

	objName := fmt.Sprintf("bilin.searchserver.SearchServantObj")
	client := bilin.NewSearchServantClient(objName, comm)

	typ := bilin.SearchType_UNKNOWN
	query := os.Args[2]
	switch os.Args[1] {
	case "user":
		typ = bilin.SearchType_USER
	case "room":
		typ = bilin.SearchType_ROOM
	case "user_room":
		typ = bilin.SearchType_USER_ROOM
	case "song":
		typ = bilin.SearchType_SONG
	}

	rsp, err := client.Search(context.Background(), &bilin.SearchReq{
		Q:   query,
		Typ: typ,
	})
	if err != nil {
		log.Fatalf("search fail: %v", err)
	}
	m := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}
	rspStr, _ := m.MarshalToString(rsp)
	log.Printf("search rsp: %s", rspStr)
}
