package handler

import (
	"bilin/protocol"
)

func retSuccess() *bilin.CommonRetInfo {
	return &bilin.CommonRetInfo{
		Ret: bilin.CommonRetInfo_RETCODE_SUCCEED,
	}
}

func retError(desc string) *bilin.CommonRetInfo {
	return &bilin.CommonRetInfo{
		Ret:  bilin.CommonRetInfo_ILLEGAL_MESSAGE,
		Desc: desc,
	}
}
