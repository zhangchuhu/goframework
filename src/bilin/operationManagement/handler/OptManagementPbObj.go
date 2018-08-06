package handler

import (
	"bilin/operationManagement/entity"
	myservice "bilin/operationManagement/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"time"
)

type OptManagementPbObj struct {
}

func NewOptManagementPbObj() *OptManagementPbObj {
	return &OptManagementPbObj{}
}

func (this *OptManagementPbObj) ActDistributionHeadgear(ctx context.Context, req *bilin.ActDistributionHeadgearRequest) (resp *bilin.ActDistributionHeadgearRespone, err error) {
	const prefix = "ActDistributionHeadgear "
	resp = &bilin.ActDistributionHeadgearRespone{Commonret: &bilin.CommonRetInfo{Ret: 0, Desc: "成功"}}

	tmEffect := time.Unix(req.Hinfo.Effecttime, 0).Format("2006-01-02 15:04:05")
	tmExpire := time.Unix(req.Hinfo.Expiretime, 0).Format("2006-01-02 15:04:05")
	headgearInfo := &entity.HeadgearInfo{int64(req.Hinfo.Uid), req.Hinfo.Headgear, tmEffect, tmExpire, req.Hinfo.Id}
	if err = myservice.MysqlReplaceVipUser(headgearInfo); err != nil {
		resp.Commonret.Ret = 1
		resp.Commonret.Desc = err.Error()
	}

	log.Debug(prefix, zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

func (this *OptManagementPbObj) BatchActDistributionHeadgear(ctx context.Context, req *bilin.BatchActDistributionHeadgearRequest) (resp *bilin.ActDistributionHeadgearRespone, err error) {
	const prefix = "BatchActDistributionHeadgear "
	resp = &bilin.ActDistributionHeadgearRespone{Commonret: &bilin.CommonRetInfo{Ret: 0, Desc: "成功"}}

	for _, item := range req.Hinfos {
		tmEffect := time.Unix(item.Effecttime, 0).Format("2006-01-02 15:04:05")
		tmExpire := time.Unix(item.Expiretime, 0).Format("2006-01-02 15:04:05")
		headgearInfo := &entity.HeadgearInfo{int64(item.Uid), item.Headgear, tmEffect, tmExpire, item.Id}
		if err = myservice.MysqlReplaceVipUser(headgearInfo); err != nil {
			resp.Commonret.Ret = 1
			resp.Commonret.Desc = err.Error()
			break
		}
	}

	log.Debug("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}

//查询某个uid用户的头像框信息
func (this *OptManagementPbObj) GetUserHeadgearInfo(ctx context.Context, req *bilin.GetUserHeadgearInfoReq) (resp *bilin.GetUserHeadgearInfoResp, err error) {
	const prefix = "GetUserHeadgearInfo "
	resp = &bilin.GetUserHeadgearInfoResp{Commonret: &bilin.CommonRetInfo{Ret: 0, Desc: "成功"}}

	sqlData, sqlErr := myservice.MysqlGetVipUser(int64(req.Uid))
	if sqlErr != nil {
		resp.Commonret.Ret = 1
		resp.Commonret.Desc = sqlErr.Error()
		return
	}
	loc, _ := time.LoadLocation("Local")
	tmEffectTime, _ := time.ParseInLocation("2006-01-02 15:04:05", sqlData.EffectTime, loc)
	tmExpireTime, _ := time.ParseInLocation("2006-01-02 15:04:05", sqlData.ExpireTime, loc)
	resp.Hinfo = &bilin.HeadgearInfo{uint64(sqlData.Uid), sqlData.Headgear, tmEffectTime.Unix(), tmExpireTime.Unix(), sqlData.Id}

	log.Debug("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp))
	return resp, nil
}
