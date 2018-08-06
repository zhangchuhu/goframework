package servant

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"os"
	"code.yy.com/yytars/goframework/jce/nodeF/taf"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
)

type NodeFHelper struct {
	comm *Communicator
	si   taf.ServerInfo
	sf   *taf.ServerF
}

func (n *NodeFHelper) SetNodeInfo(comm *Communicator, node string, app string, server string, container string) {
	n.comm = comm
	n.sf = new(taf.ServerF)
	n.sf.SetServant(comm.GetServantProxy(node))
	//comm.StringToProxy(node, n.sf)
	n.si = taf.ServerInfo{
		Application: app,
		ServerName:  server,
		Pid:         int32(os.Getpid()),
		Adapter:     "",
		ModuleType:  "taf",
		Container:   container,
	}
}

func (n *NodeFHelper) KeepAlive(adapter string) {
	appzaplog.Info("KeepAlive for adapter", zap.String("adapter", adapter))
	n.si.Adapter = adapter
	_, err := n.sf.KeepAlive(n.si)
	if err != nil {
		appzaplog.Error("keepalive fail", zap.String("adapter", adapter), zap.Error(err))
	}
}

func (n *NodeFHelper) ReportVersion(version string) {
	_, err := n.sf.ReportVersion(n.si.Application, n.si.ServerName, version)
	if err != nil {
		appzaplog.Error("report Version fail", zap.String("version", version), zap.Error(err))
	}
}
