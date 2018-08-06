package handler

import (
	"bilin/clientcenter"
	"bilin/official/dao"
	"bilin/protocol"
	"bilin/protocol/userinfocenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type LivingRecord struct {
	HOSTBilinID uint64 `json:"host_bilin_id"`
	//LivingID           int64   `json:"living_id"`
	HostNickName       string  `json:"host_nick_name"`
	LivingStartTime    string  `json:"living_start_time"`
	LivingTime         int64   `json:"living_time"` //开播时长，秒
	AudienceNum        int64   `json:"audience_num"`
	MikeUserNum        int64   `json:"mike_user_num"`
	OneMinuteOutInRate float64 `json:"one_minute_out_in_rate"`
	AverageStayTime    int64   `json:"average_stay_time"`
	RoomID             int64   `json:"room_id"`
	GiftHeartNum       int64   `json:"gift_heart_num"` //收礼心值
}

type LivingRecords struct {
	TotalPagesize int            `json:"total_pagesize"`
	Records       []LivingRecord `json:"records"`
}

func GetHostLivingRecords(c *gin.Context) *HttpError {
	appzaplog.Debug("[+]GetHostLivingRecords")
	ret := successHttp
	var (
		data *LivingRecords
	)
	for {
		hostIdInt, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetHostLivingRecords UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !accessHostRecords(c.GetInt64("uid"), hostIdInt) {
			appzaplog.Error("GetHostLivingRecords accessHostLivingRecords err", zap.Error(err))
			ret = authHttpErr
			break
		}
		data, err = hostLivingRecords(int64(hostIdInt), c)
		if err != nil {
			appzaplog.Error("GetHostLivingRecords hostLivingRecords err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
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
	appzaplog.Debug("[-]GetHostLivingRecords", zap.Any("jsonret", jsonret))
	return ret
}

func hostLivingRecords(hostid int64, c *gin.Context) (*LivingRecords, error) {
	irecords := &LivingRecords{}
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	startTimeT, err := time.Parse("20060102150405", startTime)
	if err != nil {
		return irecords, err
	}

	endTimeT, err := time.Parse("20060102150405", endTime)
	if err != nil {
		return irecords, err
	}

	page := c.Query("pageNum")
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageInt < 1 {
		return irecords, fmt.Errorf("parese pageNum err")
	}
	pagesize := c.Query("pageSize")
	pageSizeInt, err := strconv.ParseInt(pagesize, 10, 64)
	if err != nil {
		return irecords, fmt.Errorf("parese pagesize err")
	}
	info, err := dao.GetLivingRecordByHostID([]uint64{uint64(hostid)}, startTimeT.Format("2006-01-02"), endTimeT.Format("2006-01-02"))
	if err != nil {
		return irecords, err
	}
	irecords.TotalPagesize = len(info)
	end := pageInt * pageSizeInt
	start := (pageInt - 1) * pageSizeInt
	if start > int64(irecords.TotalPagesize) {
		return irecords, nil
	}
	if end > int64(irecords.TotalPagesize) {
		end = int64(irecords.TotalPagesize)
	}
	for _, v := range info[start:end] {
		irecords.Records = append(irecords.Records, LivingRecord{
			LivingStartTime:    v.LivingStartTime,
			LivingTime:         v.LivingTime,
			AudienceNum:        v.AudienceNum,
			MikeUserNum:        v.MikeUserNum,
			OneMinuteOutInRate: v.OneMinuteOutInRate,
			GiftHeartNum:       v.GiftNum,
		})
	}
	return irecords, nil
}

func GetGuildLivingRecords(c *gin.Context) *HttpError {
	ret := successHttp
	var (
		data *LivingRecords
	)
	for {
		guildIdInt, err := UInt64Param("id", c)
		if err != nil {
			appzaplog.Error("GetHostLivingRecords UInt64Param err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !owQueryGuild(c.GetInt64("uid"), guildIdInt) {
			appzaplog.Error("GetHostLivingRecords accessHostLivingRecords err", zap.Error(err))
			ret = authHttpErr
			break
		}
		data, err = guildLivingRecords(guildIdInt, c)
		if err != nil {
			appzaplog.Error("GetHostLivingRecords hostLivingRecords err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
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

func guildHostID(guildid uint64, c *gin.Context) ([]uint64, error) {
	var (
		hostids  []uint64
		blnumInt uint64
		err      error
	)

	hostBiLinNumber := c.Query("hostBiLinNumber")
	if hostBiLinNumber != "" {
		blnumInt, err = strconv.ParseUint(hostBiLinNumber, 10, 64)
		if err != nil {
			return hostids, fmt.Errorf("parese hostBiLinNumber err")
		}
	}

	if blnumInt > 0 {
		resp, err := clientcenter.UserInfoClient().BatchUserIdByBiLinId(context.TODO(), &userinfocenter.BatchUserIdByBiLinIdReq{
			Bilinid: []uint64{blnumInt},
		})
		if err != nil {
			appzaplog.Error("BatchUserIdByBiLinId err", zap.Error(err), zap.Uint64("bilinid", blnumInt))
			return hostids, fmt.Errorf("BatchUserIdByBiLinId err")
		}
		if myuid, ok := resp.Bilinid2Uid[blnumInt]; ok {
			hostids = append(hostids, myuid)
		}
	} else {
		contract, err := dao.GetContractsByGuildID(guildid)
		if err != nil {
			appzaplog.Error("GetContractsByGuildID err", zap.Error(err), zap.Uint64("guildid", guildid))
			return hostids, err
		}

		for _, v := range contract {
			if time.Now().Before(v.ContractEndTime) {
				hostids = append(hostids, v.HostUid)
			}
		}
	}
	return hostids, nil
}

func guildRoomID(guildid uint64, c *gin.Context) (rooms []uint64, err error) {
	var (
		roomidInt uint64
		guildroom *bilin.GuildRoomSResp
	)
	if roomid := c.Query("room_id"); roomid != "" {
		roomidInt, err = strconv.ParseUint(roomid, 10, 64)
		if err != nil {
			appzaplog.Error("guildRoomID ParseUint err", zap.Error(err), zap.Uint64("guildid", guildid))
			return
		}
	}
	// url里面有room_id并且大于0
	if roomidInt > 0 {
		rooms = append(rooms, roomidInt)
		return
	}

	guildroom, err = clientcenter.ConfClient().GuildRoomS(context.TODO(), &bilin.GuildRoomSReq{
		Info: &bilin.GuildRoomInfo{
			Guildid: int64(guildid),
		},
	})
	if err != nil {
		appzaplog.Error("guildLivingRecords GuildRoomS err", zap.Error(err), zap.Uint64("guildid", guildid))
		return
	}

	for _, v := range guildroom.Info {
		rooms = append(rooms, uint64(v.Roomid))
	}
	appzaplog.Debug("[-]guildRoomID", zap.Uint64("guildid", guildid), zap.Uint64s("rooms", rooms))
	return
}
func guildLivingRecords(guildid uint64, c *gin.Context) (*LivingRecords, error) {
	irecords := &LivingRecords{}
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	startTimeT, err := time.Parse("20060102150405", startTime)
	if err != nil {
		appzaplog.Error("Parse err", zap.Error(err), zap.String("startTime", startTime))
		return irecords, err
	}

	endTimeT, err := time.Parse("20060102150405", endTime)
	if err != nil {
		return irecords, err
	}

	page := c.Query("pageNum")
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageInt < 1 {
		return irecords, fmt.Errorf("parese pageNum err")
	}
	pagesize := c.Query("pageSize")
	pageSizeInt, err := strconv.ParseInt(pagesize, 10, 64)
	if err != nil {
		return irecords, fmt.Errorf("parese pagesize err")
	}

	hostids, err := guildHostID(guildid, c)
	if err != nil {
		appzaplog.Error("guildLivingRecords guildHostID err", zap.Error(err), zap.Uint64("guildid", guildid))
		return irecords, err
	}

	if len(hostids) <= 0 {
		appzaplog.Debug("guildLivingRecords empty", zap.Uint64("guildid", guildid))
		return irecords, nil
	}

	roomids, err := guildRoomID(guildid, c)
	if err != nil {
		appzaplog.Error("guildLivingRecords guildRoomID err", zap.Error(err), zap.Uint64("guildid", guildid))
		return irecords, err
	}
	if len(roomids) <= 0 {
		appzaplog.Debug("guildLivingRecords roomid empty", zap.Uint64("guildid", guildid))
		return irecords, nil
	}

	info, err := dao.GetLivingRecordByHostIDAndRoomID(hostids, roomids, startTimeT.Format("2006-01-02"), endTimeT.Format("2006-01-02"))
	if err != nil {
		appzaplog.Error("GetLivingRecordByHostIDAndRoomID err", zap.Error(err), zap.Uint64("guildid", guildid))
		return irecords, err
	}
	irecords.TotalPagesize = len(info)
	end := pageInt * pageSizeInt
	start := (pageInt - 1) * pageSizeInt
	if start > int64(irecords.TotalPagesize) {
		return irecords, nil
	}
	if end > int64(irecords.TotalPagesize) {
		end = int64(irecords.TotalPagesize)
	}
	uinfo, err := clientcenter.TakeUserInfo(hostids)
	if err != nil {
		appzaplog.Error("TakeUserInfo err", zap.Error(err))
		return irecords, err
	}

	bilinid, err := clientcenter.UserInfoClient().BatchUserBiLinId(context.TODO(), &userinfocenter.BatchUserBiLinIdReq{
		Uid: hostids,
	})
	if err != nil {
		appzaplog.Error("BatchUserBiLinId err", zap.Error(err))
		return irecords, err
	}

	for _, v := range info[start:end] {
		rec := LivingRecord{
			LivingStartTime:    v.LivingStartTime,
			LivingTime:         v.LivingTime,
			AudienceNum:        v.AudienceNum,
			MikeUserNum:        v.MikeUserNum,
			OneMinuteOutInRate: v.OneMinuteOutInRate,
			RoomID:             v.RoomID,
			GiftHeartNum:       v.GiftNum,
		}
		uhostid, err := strconv.ParseUint(v.HOSTID, 10, 64)
		if err != nil {
			appzaplog.Error("ParseUint err", zap.Error(err), zap.String("hostid", v.HOSTID))
			continue
		}
		if bilinid, ok := bilinid.Uid2Bilinid[uhostid]; ok {
			rec.HOSTBilinID = bilinid
		}
		if nick, ok := uinfo[uhostid]; ok {
			rec.HostNickName = nick.NickName
		}
		irecords.Records = append(irecords.Records, rec)
	}
	return irecords, nil
}
