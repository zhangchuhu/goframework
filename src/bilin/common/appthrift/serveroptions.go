// @author kordenlu
// @创建时间 2017/12/01 11:04
// 功能描述:

package appthrift

import "git.apache.org/thrift.git/lib/go/thrift"

type ThriftServerOption func(*AppThriftServer)

func ThriftHost(host string) ThriftServerOption {
	return func(this *AppThriftServer) {
		this.listerHost = host
	}
}

func TransportFactory(tf thrift.TTransportFactory) ThriftServerOption {
	return func(this *AppThriftServer) {
		this.transportFactory = tf
	}
}

func ProtocolFactory(pf thrift.TProtocolFactory) ThriftServerOption {
	return func(this *AppThriftServer) {
		this.protocolFactory = pf
	}
}
