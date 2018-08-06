package handler

import (
	"bilin/confinfocenter/dao"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"sort"
	"time"
)

type ConfInfoServantObj struct {
}

func NewConfInfoServantObj() *ConfInfoServantObj {
	return &ConfInfoServantObj{}
}

type LivingcategorySlice []*bilin.LivingCategoryInfo

func (a LivingcategorySlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a LivingcategorySlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a LivingcategorySlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Sort < a[i].Sort
}

// 获取直播间分类信息配置
func (this *ConfInfoServantObj) LivingCategorys(ctx context.Context, r *bilin.LivingCategorysReq) (*bilin.LivingCategorysResp, error) {
	code := int64(0)
	defer func(now time.Time) {
		httpmetrics.DefReport("LivingCategorys", code, now, httpmetrics.DefaultSuccessFun)
	}(time.Now())
	category, err := dao.GetLivingCategorys()
	if err != nil {
		code = GetLivingCategorysFailed
		appzaplog.Error("[+]LivingCategorys GetLivingCategorys failed", zap.Error(err))
		return nil, err
	}
	resp := &bilin.LivingCategorysResp{}
	for _, v := range category {
		resp.Livingcategorys = append(resp.Livingcategorys, &bilin.LivingCategoryInfo{
			Typeid:         int64(v.TYPE_ID),
			Name:           v.TYPE_NAME,
			Color:          v.FONT_COLOR,
			Backgroudimage: v.BACKGROUND_IMAGE,
			Sort:           int64(v.SORT),
		})
	}
	sort.Sort(LivingcategorySlice(resp.Livingcategorys))
	return resp, nil
}
