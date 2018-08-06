package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type HttpRetComm struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type PageInfo struct {
	pageNum  int
	pageSize int
}

func pageInfo(c *gin.Context) (PageInfo, error) {
	info := PageInfo{}
	pageNum := c.DefaultQuery("page_num", "0")
	pageNumInt, err := strconv.ParseInt(pageNum, 10, 64)
	if err != nil {
		return info, err
	}
	info.pageNum = int(pageNumInt)

	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	if err != nil {
		return info, err
	}
	info.pageSize = int(pageSize)
	return info, nil
}

func (p *PageInfo) StartEnd(totallen int) (start, end int, err error) {
	if p.pageNum < 0 {
		err = fmt.Errorf("pageNum too small:%d", p.pageNum)
		return
	}
	if p.pageNum == 0 {
		end = totallen
		return
	}
	start = p.pageSize * (p.pageNum - 1)
	end = p.pageSize * p.pageNum
	if end > totallen {
		end = totallen
	}
	return
}
