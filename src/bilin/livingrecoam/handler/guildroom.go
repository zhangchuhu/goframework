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
)

type GuildRoom struct {
	ID      int64 `json:"id"`
	GuildID int64 `json:"guild_id"`
	RoomID  int64 `json:"room_id"`
}

type GuildRoomList struct {
	TotalPageSize int         `json:"total_page_size"`
	Items         []GuildRoom `json:"items"`
}

func GetGuildRoom(c *gin.Context) {
	appzaplog.Debug("[+]GetGuildRoom")
	code := http.StatusOK
	guildRec := &GuildRoomList{}

	for {
		pageinfo, err := pageInfo(c)
		if err != nil {
			appzaplog.Error("GetGuildRoom pageInfo err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		info, err := clientcenter.GuildTarsClient().RGuildRoom(context.TODO(), &bilin.RGuildRoomReq{})
		if err != nil {
			appzaplog.Error("GetGuildRoom GuildRoomS err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		guildRec.TotalPageSize = len(info.Info)
		for _, v := range info.Info {
			guildRec.Items = append(guildRec.Items, GuildRoom{
				RoomID:  v.Roomid,
				GuildID: v.Guildid,
				ID:      v.Id,
			})
		}

		if len(guildRec.Items) > 0 {
			start, end, err := pageinfo.StartEnd(len(guildRec.Items))
			if err != nil {
				appzaplog.Error("GetGuildRoom StartEnd err", zap.Error(err))
				code = http.StatusBadRequest
				break
			}
			guildRec.Items = guildRec.Items[start:end]
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]GetGuildRoom", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func CreateGuildRoom(c *gin.Context) {
	appzaplog.Debug("[+]CreateGuildRoom")
	code := http.StatusOK
	guildRec := &GuildRoom{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("CreateGuildRoom BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("CreateGuildRoom", zap.Any("req", guildRec))

		_, err := clientcenter.GuildTarsClient().CGuildRoom(context.TODO(), &bilin.CGuildRoomReq{
			Info: &bilin.GuildRoom{
				Id:      guildRec.ID,
				Roomid:  guildRec.RoomID,
				Guildid: guildRec.GuildID,
			},
		})
		if err != nil {
			appzaplog.Error("CreateGuildRoom CreateGuildRoom err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]CreateGuildRoom", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func UpdateGuildRoom(c *gin.Context) {
	appzaplog.Debug("[+]UpdateGuildRoom")
	code := http.StatusOK
	guildRec := &GuildRoom{}

	for {
		if err := c.BindJSON(guildRec); err != nil {
			appzaplog.Error("UpdateGuildRoom BindJSON err", zap.Error(err))
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("UpdateGuildRoom", zap.Any("req", guildRec))

		_, err := clientcenter.GuildTarsClient().UGuildRoom(context.TODO(), &bilin.UGuildRoomReq{
			Info: &bilin.GuildRoom{
				Id:      guildRec.ID,
				Roomid:  guildRec.RoomID,
				Guildid: guildRec.GuildID,
			},
		})
		if err != nil {
			appzaplog.Error("UpdateGuildRoom BindJSON err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
		Data: guildRec,
	}

	appzaplog.Debug("[-]UpdateGuildRoom", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}

func DelGuildRoom(c *gin.Context) {
	appzaplog.Debug("[+]DelGuildRoom")
	code := http.StatusOK

	for {
		roomid := c.Param("id")
		if roomid == "" {
			appzaplog.Error("DelGuildRoom empty roomid err")
			code = http.StatusBadRequest
			break
		}
		appzaplog.Debug("DelGuildRoom", zap.String("roomid", roomid))
		roomidU64, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("DelGuildRoom ParseUint err")
			code = http.StatusBadRequest
			break
		}
		_, err = clientcenter.GuildTarsClient().DGuildRoom(context.TODO(), &bilin.DGuildRoomReq{
			Info: &bilin.GuildRoom{
				Id: roomidU64,
			},
		})
		if err != nil {
			appzaplog.Error("DelGuildRoom DelCategoryGuildRec err", zap.Error(err))
			code = http.StatusInternalServerError
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: code,
	}

	appzaplog.Debug("[-]DelGuildRoom", zap.Any("resp", jsonret))
	c.JSON(200, jsonret)
}
