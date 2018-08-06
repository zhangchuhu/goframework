package nsqwrapper

import (
	"bilin/userinfocenter/handler"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"github.com/nsqio/go-nsq"
)

const (
	updateUserInfoTopic = "updatebilinuserinfo"
)

func InitNsq(nsqlookupaddrs []string) error {
	cfg := nsq.NewConfig()
	//cfg.Set("local_addr", "127.0.0.1:4150")
	cfg.UserAgent = fmt.Sprintf("nsq_pubsub/%s go-nsq/%s", "userinfo", nsq.VERSION)
	cfg.MaxInFlight = 5
	r, err := nsq.NewConsumer(updateUserInfoTopic, "nick", cfg)
	if err != nil {
		return err
	}
	r.AddHandler(&handler.UserInfoEventHandler{})
	if err := ConnectToNSQAndLookupd(r, nsqlookupaddrs); err != nil {
		return err
	}
	go func() {
		for range r.StopChan {
			appzaplog.Warn("NickConsumer stoped")
		}
	}()
	return nil
}

func ConnectToNSQAndLookupd(r *nsq.Consumer, lookupd []string) error {
	//for _, addrString := range nsqAddrs {
	//	err := r.ConnectToNSQD(addrString)
	//	if err != nil {
	//		return err
	//	}
	//}

	for _, addrString := range lookupd {
		appzaplog.Debug("lookupd addr %s", zap.String("lookupaddr", addrString))
		err := r.ConnectToNSQLookupd(addrString)
		if err != nil {
			return err
		}
	}

	return nil
}
