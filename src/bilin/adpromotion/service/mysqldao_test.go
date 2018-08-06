package service

import (
	"bilin/adpromotion/config"
	"testing"
)

var (
	appconfig = &config.AppConfig{
		MysqlAddr: "adpromotion:xptKokbu7@tcp(221.228.110.28:6315)/bilin_adpromotion?charset=utf8",
	}
)

func init() {
	config.SetTestAppConfig(appconfig)
	MysqlInit()
}

func TestMysqlSelectClickInfoByImei(t *testing.T) {
	_, err := MysqlSelectClickInfoByImei("e807f1fcf82d132f9bb018ca6738a19f")
	if err != nil {
		t.Error(err)
		return
	}

	//查找db，如果命中，则上报数据到360平台
	clickInfo, errSelect := MysqlSelectClickInfoByImei("e807f1fcf82d132f9bb018ca6738a19f")
	if errSelect != nil {
		t.Error(err)
	}

	// report qihu360, http get
	resp, err := ReportQihu360(clickInfo)
	if err != nil {
		t.Error(err)
	}

	UpdateQihu360CallBackResult(clickInfo, resp)

	t.Log(resp)
}
