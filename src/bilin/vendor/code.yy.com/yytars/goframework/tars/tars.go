// @author kordenlu
// @创建时间 2018/03/09 17:10
// 功能描述: 框架对外接口统一放置

package tars

import (
	"code.yy.com/yytars/goframework/tars/servant"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"context"
)

// AddServant only used for tars header,idl plugin use tarsheaderpbbody2go
func AddServant(v servant.Dispatcher, f interface{}, objname string) error{
	return servant.AddServant(v,f,objname)
}

//
func Run()  {
	servant.Run()
}

// 根据文件名获取应用配置,如果不存在会尝试拉取一次
func ReadConf(filename string) ([]byte,error) {
	return servant.ReadConf(filename)
}

//订阅文件配置变更
func SubConfigPush() <-chan *pbtaf.ConfigPushNotice{
	return servant.SubConfigPush()
}

// for client
func NewCommunicator() *servant.Communicator{
	return servant.NewPbCommunicator()
}

// 客户端调用时，可以设置用户的透传信息到ctxmap
// server端会通过FromOutgoingContext获取到ctxmap
func NewOutgoingContext(ctx context.Context, ctxmap map[string]string) context.Context {
	return servant.NewOutgoingContext(ctx,ctxmap)
}

// server端通过这个函数获取client在NewOutgoingContext中设置的信息
func FromOutgoingContext(ctx context.Context) (md map[string]string, ok bool) {
	return servant.FromOutgoingContext(ctx)
}