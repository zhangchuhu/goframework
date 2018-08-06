package servant

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/jce/notify/taf"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
)

type NotifyHelper struct {
	comm *Communicator
	tn   *taf.Notify
	tm   taf.ReportInfo
}

var defaultNotifyHelper *NotifyHelper

func initNotify(comm *Communicator,srvconfig *serverConfig) error  {
	if comm == nil || srvconfig == nil{
		return NilParamsErr
	}
	tn := &taf.Notify{}
	tn.SetServant(comm.GetServantProxy(srvconfig.notify))
	defaultNotifyHelper = &NotifyHelper{
		comm:comm,
		tn :tn,
		//todo params
		tm : taf.ReportInfo{
			SApp:srvconfig.App,
			SServer:srvconfig.Server,
			SContainer:srvconfig.Container,
		},
	}
	return nil
}

func (n *NotifyHelper) ReportNotifyInfo(info string) {
	tm := n.tm
	tm.SMessage = info
	appzaplog.Debug("ReportNotifyInfo", zap.Any("ReportInfo", tm))
	n.tn.ReportNotifyInfo(tm)
}

func reportNotifyInfo(info string) {
	if defaultNotifyHelper == nil{
		appzaplog.Error("notify client not init")
		return
	}
	defer func() {
		if err := recover(); err != nil {
			appzaplog.Debug("recover", zap.Any("err", err))
		}
	}()
	defaultNotifyHelper.ReportNotifyInfo(info)
}
