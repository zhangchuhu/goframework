package handler

import (
	"bilin/confinfocenter/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"time"
)

// 审核过滤关键字
func (this *ConfInfoServantObj) AppleAuditWords(ctx context.Context, r *bilin.AppleAuditWordsReq) (*bilin.AppleAuditWordsResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("AppleAuditWords", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	words, err := dao.GetAuditWorld()
	if err != nil {
		code = GetAuditWorldFailed
		return nil, err
	}
	ret := &bilin.AppleAuditWordsResp{
		Auditwords: words,
	}
	return ret, nil
}
