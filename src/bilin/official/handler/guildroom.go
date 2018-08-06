/*
工会房间
*/
package handler

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

type GuildRoom struct {
	GuildId int64   `json:"guild_id"`
	RoomIds []int64 `json:"room_ids"`
}

func GetGuildRoom(c *gin.Context) *HttpError {
	var (
		ret  = successHttp
		data = &GuildRoom{}
	)

	for {
		cookieuid := c.GetInt64("uid")
		id := c.Param("id")
		guildIdInt, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			appzaplog.Error("GetGuildRoom ParseUint err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}

		if !owQueryGuild(cookieuid, guildIdInt) {
			appzaplog.Warn("GetGuildRoom not ow", zap.Uint64("guildid", guildIdInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}

		// 获取签约信息
		resp, err := clientcenter.ConfClient().GuildRoomS(context.TODO(), &bilin.GuildRoomSReq{
			Info: &bilin.GuildRoomInfo{
				Guildid: int64(guildIdInt),
			},
		})
		if err != nil {
			appzaplog.Error("GetGuildRoom Contract err", zap.Uint64("guildid", guildIdInt), zap.Error(err))
			ret = daoGetHttpErr
			break
		}
		data.GuildId = int64(guildIdInt)
		for _, v := range resp.Info {
			data.RoomIds = append(data.RoomIds, v.Roomid)
		}
		appzaplog.Debug("[+]GetHostGuild", zap.String("id", id))

		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 {
		jsonret.Data = data
	}
	c.JSON(200, jsonret)
	return ret
}
