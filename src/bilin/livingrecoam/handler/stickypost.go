package handler

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type StickyPost struct {
	ID        int64  `json:"id,omitempty"`
	Roomid    uint64 `json:"roomid,omitempty"`
	Typeid    int64  `json:"typeid,omitempty"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Weight    int64  `json:"weight"`
}

type StickyPostList struct {
	TotalPageSize int          `json:"total_page_size"`
	Items         []StickyPost `json:"items"`
}

const TimeLayoutOther string = "2006-01-02 15:04:05"

func GetStickyPost(c *gin.Context) {
	appzaplog.Debug("[+]GetStickyPost")
	code := http.StatusOK
	stickyPostRec := &StickyPostList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetStickyPost pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.ConfClient().AvailableCategoryStickie(context.TODO(), &bilin.AvailableCategoryStickieReq{})
		if err != nil {
			appzaplog.Error("GetStickyPost err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		stickyPostRec.TotalPageSize = len(info.Infos)
		for _, v := range info.Infos {
			stickyPostRec.Items = append(stickyPostRec.Items, StickyPost{
				ID:        v.Id,
				Roomid:    v.Roomid,
				Typeid:    v.Typeid,
				StartTime: time.Unix(v.Starttime, 0).Format(time.RFC3339Nano),
				EndTime:   time.Unix(v.Endtime, 0).Format(time.RFC3339Nano),
				Weight:    v.Sort,
			})
		}
		if len(stickyPostRec.Items) > 0 {
			//hostRec.TotalPageSize = len(hostRec.Items)/pageinfo.pageSize + 1
			start, end, err := pageinfo.StartEnd(len(stickyPostRec.Items))
			if err != nil {
				appzaplog.Error("GetStickyPost StartEnd err", zap.Error(err))
				code = http.StatusBadRequest
				break
			}
			stickyPostRec.Items = stickyPostRec.Items[start:end]
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: stickyPostRec,
	}

	appzaplog.Debug("[-]GetStickyPost", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func AddStickyPost(c *gin.Context) {
	appzaplog.Debug("[+]AddStickyPost")
	code := http.StatusOK
	guildRec := &StickyPost{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("AddStickyPost BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("AddStickyPost", zap.Any("req", guildRec))

		startTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.StartTime, time.Local)
		if err != nil {
			appzaplog.Error("AddStickyPost ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}

		endTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.EndTime, time.Local)
		if err != nil {
			appzaplog.Error("AddStickyPost ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}

		_, err = clientcenter.ConfClient().CreateCategoryStickie(context.TODO(), &bilin.CreateCategoryStickieReq{
			Info: &bilin.CategoryStickieInfo{
				Roomid:    guildRec.Roomid,
				Typeid:    guildRec.Typeid,
				Starttime: startTime.Unix(),
				Endtime:   endTime.Unix(),
				Sort:      guildRec.Weight,
			},
		})
		if err != nil {
			appzaplog.Error("AddStickyPost CreateCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]AddStickyPost", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelStickyPost(c *gin.Context) {
	appzaplog.Debug("[+]DelStickyPost")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelStickyPost empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelStickyPost", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelStickyPost ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ConfClient().DelCategoryStickie(context.TODO(), &bilin.DelCategoryStickieReq{
			Info: &bilin.CategoryStickieInfo{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelStickyPost DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelGuildRec", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateStickyPost(c *gin.Context) {
	appzaplog.Debug("[+]UpdateStickyPost")
	code := http.StatusOK
	guildRec := &StickyPost{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateStickyPost BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateStickyPost", zap.Any("req", guildRec))
		startTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.StartTime, time.Local)
		if err != nil {
			appzaplog.Error("AddStickyPost ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}

		endTime, err := time.ParseInLocation(time.RFC3339Nano, guildRec.EndTime, time.Local)
		if err != nil {
			appzaplog.Error("AddStickyPost ParseInLocation err", zap.Error(err), zap.Any("req", guildRec))
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.ConfClient().UpdateCategoryStickie(context.TODO(), &bilin.UpdateCategoryStickieReq{
			Info: &bilin.CategoryStickieInfo{
				Id:        guildRec.ID,
				Typeid:    guildRec.Typeid,
				Sort:      guildRec.Weight,
				Roomid:    guildRec.Roomid,
				Starttime: startTime.Unix(),
				Endtime:   endTime.Unix(),
			},
		})
		if err != nil {
			appzaplog.Error("UpdateStickyPost BindJSON err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateStickyPost", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
