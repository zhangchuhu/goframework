package sensitive

import (
	"bilin/clientcenter"
	"bilin/protocol"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	st "github.com/importcjj/sensitive"
)

var filter *st.Filter

func LoadSensiTive() error {
	resp, err := clientcenter.ConfClient().AppleAuditWords(context.TODO(), &bilin.AppleAuditWordsReq{})
	if err != nil {
		appzaplog.Error("AppleAuditWords failed", zap.Error(err))
		return err
	}
	filter = st.New()
	appzaplog.Info("LoadSensiTive", zap.Strings("sensitiveword", resp.Auditwords))
	filter.AddWord(resp.Auditwords...)
	return nil
}

func SensiTiveWord(word string) bool {
	if filter == nil {
		appzaplog.Warn("filter not init")
		return false
	}
	ret, _ := filter.FindIn(word)
	return ret
}
