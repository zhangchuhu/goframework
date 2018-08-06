// @author kordenlu
// @创建时间 2018/02/13 11:09
// 功能描述:

package servant

import "errors"

//
var (
	// comm err
	NilParamsErr = errors.New("nil params")

	// objectproxy error
	NoAdapterErr  = errors.New("no adapter proxy selected")
	OverloadErr   = errors.New("invoke queue is full")
	ReqTimeoutErr = errors.New("request timeout")
	SendErr       = errors.New("send failed")

	// server dispatch error
	//HandlerPanicErr    = errors.New("handler panic")
	//RPCCallDisabledErr = errors.New("rpc call disabled")
)

type TarError struct {
	Code int64
	Desc string
}

func (err TarError)Error() string{
	return err.Desc
}

const (
	TarsAppErrBegin = 1
	TarsAppErrEnd = 10000

	TarsServerErrBegin = 10001
	TarsServerErrEnd = 20000

	TarsClientErrBegin = 20001
	TarsClientErrEnd = 30000
)

var(
	// 0 成功
	TarSuccess = TarError{
		Code:0,
		Desc:"success",
	}

	// code 1 - 10000 用于业务error定义

	// code 10001 - 20000用于tars框架server端err定义
	HandlerPanicTarErr    = TarError{
		Code:TarsServerErrBegin+1,
		Desc:"handler panic",
	}
	RPCCallDisabledTarErr = TarError{
		Code:TarsServerErrBegin+2,
		Desc:"rpc call disabled",
	}

	// code 20001 - 30000 reserved by tars client err
	NoAdapterTarErr  = TarError{
		Code:TarsClientErrBegin+1,
		Desc:"no adapter proxy selected",
	}
	OverloadTarErr   = TarError{
		Code:TarsClientErrBegin+2,
		Desc:"invoke queue is full",
	}
	ReqTimeoutTarErr = TarError{
		Code:TarsClientErrBegin+3,
		Desc:"request timeout",
	}

	SendTarErr       = TarError{
		Code:TarsClientErrBegin+4,
		Desc:"send failed",
	}
	// 剩余部分保留
)