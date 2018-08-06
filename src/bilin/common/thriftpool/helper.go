package thriftpool

import (
	"fmt"
	"time"

	"bilin/bcserver/bccommon"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
)

// Run accept the generic service client
type Run func(interface{}) error

func Invoke(prefix string, pool Pool, run Run) (err error) {
	var metrics_ret int64 = 0
	defer func(now time.Time) {
		httpmetrics.DefReport(prefix, metrics_ret, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	conn, err := pool.Get()
	if err != nil {
		log.Error(fmt.Sprintf("%q can not get thrift conn from pool: %v", prefix, err))
		metrics_ret = -1
		return
	}
	defer conn.Close()
	err = run(conn.Client)
	if err != nil {
		conn.Conn.Close()                                               // close the socket that failed
		if conn.Conn, err = pool.(*channelPool).factory(); err != nil { // reconnect the socket
			log.Error(fmt.Sprintf("%q reconnect failed, server down? %v", prefix, err))
			conn.MarkUnusable()
			metrics_ret = -1
			return
		}
		err = run(conn.Client) // retry on the newly connected socket
		if err != nil {
			log.Error(fmt.Sprintf("%q failed after reconnect, fatal! %v", prefix, err))
			conn.MarkUnusable()
			metrics_ret = -1
			return
		}
	}
	return
}

func Ping(prefix string, pool Pool, run Run, interval time.Duration) {
	var count int
	for {
		count = 0
		n := pool.Len()
		for i := 0; i < n; i++ {
			conn, err := pool.Get()
			if err != nil {
				break
			}
			err = run(conn.Client)
			if err != nil {
				count++
				conn.MarkUnusable()
			}
			conn.Close()
		}
		if count > 0 {
			log.Info(fmt.Sprintf("%q removed %d stale thrift connection(s) out of %d in the pool", prefix, count, n))
		}
		time.Sleep(interval)
	}
}
