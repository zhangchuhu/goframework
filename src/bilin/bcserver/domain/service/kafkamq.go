package service

import (
	"bilin/bcserver/config"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/Shopify/sarama"
)

var (
	producer sarama.SyncProducer
)

func KafkaProducerInit() {
	const prefix = "KafkaProducerInit "

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	var err error
	producer, err = sarama.NewSyncProducer(config.GetAppConfig().KafkaAddr, kafkaConfig)
	if err != nil {
		log.Error(prefix+"NewSyncProducer Fail", zap.Any("err", err))
		panic(err)
	}

	log.Info(prefix + "NewSyncProducer success")
}

func KafkaProduceMessage(msg string) (err error) {
	const prefix = "KafkaProduceMessage "

	sendMsg := &sarama.ProducerMessage{
		Topic:     config.GetAppConfig().KafkaTopic,
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}
	sendMsg.Value = sarama.ByteEncoder(msg)
	paritition, offset, err := producer.SendMessage(sendMsg)
	if err != nil {
		log.Error(prefix+"Send Message Failed", zap.Any("err", err))
		return
	}

	log.Info(prefix+"end", zap.Any("msg", msg), zap.Any("paritition", paritition), zap.Any("offset", offset))
	return
}
