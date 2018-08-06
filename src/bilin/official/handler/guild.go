package handler

import (
	"bilin/official/dao"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/gin-gonic/gin"
	"strconv"
)

type GuildInfo struct {
	GuildID   uint64 `json:"guild_id"`
	OW        uint64 `json:"ow"`
	Title     string `json:"title"`
	Mobile    string `json:"mobile"`
	Desc      string `json:"desc"`
	GuildLogo string `json:"guild_logo"`
}

//const corsACAO = "http://pgbilin.yy.com"
//const corsACAO = "http://172.27.142.9"

func GetHostGuild(c *gin.Context) *HttpError {
	var (
		ret      = successHttp
		myguild  *dao.Guild
		contract *dao.Contract
	)

	for {
		cookieuid := c.GetInt64("uid")
		id := c.Param("id")
		idInt, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			appzaplog.Error("GetHostGuild ParseUint err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}

		if !cookieUserEqReqUser(cookieuid, idInt) {
			appzaplog.Warn("GetHostGuild not host", zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}

		// 获取签约信息
		if contract, err = dao.GetContractByHostUid(idInt); err != nil {
			appzaplog.Error("GetHostGuild Contract err", zap.Uint64("uid", idInt), zap.Error(err))
			ret = daoGetHttpErr
			break
		}
		if contract == nil {
			appzaplog.Debug("no contract", zap.Int64("uid", cookieuid), zap.Uint64("id", idInt))
			break
		}
		appzaplog.Debug("[+]GetHostGuild", zap.String("id", id))

		if myguild, err = dao.GetByGuildID(contract.GuildID); err != nil {
			appzaplog.Error("GetHostGuild Guild err", zap.Uint64("uid", idInt), zap.Error(err))
			ret = daoGetHttpErr
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 && myguild != nil {
		jsonret.Data = &GuildInfo{
			GuildID:   uint64(myguild.ID),
			OW:        myguild.OW,
			Title:     myguild.Title,
			Mobile:    myguild.Mobile,
			Desc:      myguild.Describle,
			GuildLogo: myguild.GuildLogo,
		}
	}
	c.JSON(200, jsonret)
	return ret
}

func GetGuildByOwID(c *gin.Context) *HttpError {

	var (
		ret     = successHttp
		myguild = &dao.Guild{}
	)
	for {
		cookieuid := c.GetInt64("uid")
		id := c.Param("id")
		idInt, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			appzaplog.Error("GetGuildByOwID ParseUint err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}
		if !cookieUserEqReqUser(cookieuid, idInt) {
			appzaplog.Warn("GetGuildByOwID not OW", zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}
		if myguild, err = dao.GetGuildByOW(idInt); err != nil {
			appzaplog.Error("GetGuildByOwID GetGuildByOW err", zap.Uint64("uid", idInt), zap.Error(err))
			ret = daoGetHttpErr
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 && myguild != nil {
		jsonret.Data = &GuildInfo{
			GuildID:   uint64(myguild.ID),
			OW:        myguild.OW,
			Title:     myguild.Title,
			Mobile:    myguild.Mobile,
			Desc:      myguild.Describle,
			GuildLogo: myguild.GuildLogo,
		}
	}
	c.JSON(200, jsonret)
	return ret
}

func UpdateGuildByOwID(c *gin.Context) *HttpError {
	var (
		ret         = successHttp
		updateguild *dao.Guild
	)
	for {
		cookieuid := c.GetInt64("uid")
		id := c.Param("id")
		idInt, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			appzaplog.Error("UpdateGuildByOwID ParseUint err", zap.Error(err))
			ret = parseURLHttpErr
			break
		}

		if !cookieUserEqReqUser(cookieuid, idInt) {
			appzaplog.Warn("UpdateGuildByOwID not OW", zap.Uint64("requid", idInt), zap.Int64("cookieuid", cookieuid))
			ret = authHttpErr
			break
		}

		updateguild = takePutGuild(c)
		if err := updateguild.UpdateByOw(idInt); err != nil {
			appzaplog.Error("UpdateGuildByOwID err", zap.Error(err))
			ret = daoPutHttpErr
			break
		}
		break
	}

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	if ret.code == 0 {
		jsonret.Data = &GuildInfo{
			Title:     updateguild.Title,
			Mobile:    updateguild.Mobile,
			Desc:      updateguild.Describle,
			GuildLogo: updateguild.GuildLogo,
		}
	}
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	c.JSON(201, jsonret)
	return ret
}

func OptionGuildByOwID(c *gin.Context) {
	appzaplog.Debug("[+]OptionGuildByOwID Option")
	var (
		ret = successHttp
	)

	jsonret := &HttpRetComm{
		Code: ret.code,
		Desc: ret.desc,
	}
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")

	c.JSON(204, jsonret)
	appzaplog.Debug("[-]OptionGuildByOwID Option", zap.Any("resp", jsonret))
}

// title=""&mobile=""&desc=""&guild_logo=""
func takePutGuild(c *gin.Context) *dao.Guild {
	appzaplog.Debug("[+]takePutGuild", zap.Int64("uid", c.GetInt64("uid")))
	ret := &dao.Guild{}
	if title, ok := c.GetPostForm("title"); ok {
		ret.Title = title
	}
	if mobile, ok := c.GetPostForm("mobile"); ok {
		ret.Mobile = mobile
	}

	if desc, ok := c.GetPostForm("desc"); ok {
		ret.Describle = desc
	}

	if guild_log, ok := c.GetPostForm("guild_logo"); ok {
		ret.GuildLogo = guild_log
	}

	appzaplog.Debug("[-]takePutGuild", zap.Any("resp", ret))
	return ret
}
