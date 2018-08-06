// @author kordenlu
// @创建时间 2018/03/19 10:31
// 功能描述:

package sd

import (
	"code.yy.com/yytars/goframework/jce/servant/taf"
	"github.com/juju/ratelimit"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"errors"
)

var (
	RateLimiterErr = errors.New("rate limiter triggered")
)

type SDHelper interface {
	FindObjectByIdInSameGroup(id string,activeEp *[]taf.EndpointF,inactiveEp *[]taf.EndpointF,_opt ...map[string]string )(_ret int32,_err error)
}

type QueryFHelper struct {
	qratelimiter *ratelimit.Bucket
	q            *taf.QueryF
}

func NewQueryFHelper(qratelimiter *ratelimit.Bucket, q *taf.QueryF) SDHelper {
	return &QueryFHelper{
		qratelimiter:qratelimiter,
		q:q,
	}
}

func (this *QueryFHelper)FindObjectByIdInSameGroup(id string,activeEp *[]taf.EndpointF,inactiveEp *[]taf.EndpointF,_opt ...map[string]string )(_ret int32,_err error){
	if this.qratelimiter != nil && this.qratelimiter.TakeAvailable(1) == 0{
		appzaplog.Warn("FindObjectByIdInSameGroup rate limite triggered")
		return 0,RateLimiterErr
	}
	return this.q.FindObjectByIdInSameGroup(id,activeEp,inactiveEp)
}