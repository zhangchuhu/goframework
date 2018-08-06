package service

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster" //support automatic consumer-group rebalancing and offset tracking
	"strings"
	"time"

	"bilin/adpromotion/config"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func handleKafkaMsg(message []byte) {
	const prefix = "handleKafkaMsg "

	var result map[string][]interface{}
	err := json.Unmarshal(message, &result)
	if err != nil {
		log.Error(prefix+"json.Unmarshal", zap.Any("err", err), zap.Any("message", string(message)))
		return
	}

	datas := result["datas"]
	for _, item := range datas {
		log.Debug(prefix, zap.Any("item", item))

		imei := item.(map[string]interface{})["imei"]
		if imei == nil {
			log.Warn(prefix, zap.Any("imei", imei))
			continue
		}

		From := item.(map[string]interface{})["from"]
		if From == nil {
			log.Warn(prefix, zap.Any("From", From))
			continue
		}

		if From.(string) != "360" {
			continue
		}

		Imei_md5 := GetMd5String(imei.(string))
		log.Debug(prefix, zap.Any("from", From.(string)), zap.Any("Imei_md5", Imei_md5))

		//查找db，如果命中，则上报数据到360平台
		clickInfo, errSelect := MysqlSelectClickInfoByImei(Imei_md5)
		if errSelect != nil {
			log.Warn(prefix+"MysqlSelectClickInfoByImei", zap.Any("errSelect", errSelect))
			continue
		}

		// report qihu360, http get
		resp, err := ReportQihu360(clickInfo)
		if err != nil {
			log.Error(prefix, zap.Any("ReportQihu360", err))
			continue
		}

		UpdateQihu360CallBackResult(clickInfo, resp)
	}

}

func ConsumerKafka() {
	const prefix = "ConsumerKafka "
	log.Debug(prefix + "begin")

	groupID := "group-1"
	kafkaConfig := cluster.NewConfig()
	kafkaConfig.Group.Return.Notifications = true
	kafkaConfig.Consumer.Offsets.CommitInterval = 1 * time.Second
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始

	c, err := cluster.NewConsumer(config.GetAppConfig().KafkaAddr, groupID, strings.Split(config.GetAppConfig().Topic, ","), kafkaConfig)
	if err != nil {
		log.Error(prefix, zap.Any("Failed open consumer", err))
		return
	}
	defer c.Close()
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		noti := c.Notifications()
		for {
			select {
			case err := <-errors:
				log.Error(prefix, zap.Any("Failed Notifications", err))
			case <-noti:
			}
		}
	}(c)

	for msg := range c.Messages() {
		log.Debug(prefix, zap.Any("recv Messages", string(msg.Value)))
		c.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset

		handleKafkaMsg(msg.Value)
	}
}
