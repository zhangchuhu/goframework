package service

import (
	"bilin/relationlist/config"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"encoding/json"
	"github.com/streadway/amqp"
)

const (
	ExchangeName = "to.pub.external"
	QueueName    = "Bilin.props.queue"
	RoutingKey   = "props.use.Bilin"
)

//用户数据
type PropsInfo struct {
	TotalAmount uint64 `json:"totalAmount"`
}

//榜单数据
type RabbitMqMessage struct {
	ToUid     uint64      `json:"recvUid"`
	FromUid   uint64      `json:"uid"`
	PropsList []PropsInfo `json:"usedInfo"`
}

func handleRabbitMqMsg(message []byte) {
	const prefix = "handleRabbitMqMsg "

	recvMsg := &RabbitMqMessage{}
	if err := json.Unmarshal(message, recvMsg); err != nil {
		log.Warn(prefix+"Unmarshal failed", zap.Any("err", err))
		return
	}

	var totalAmount uint64 = 0
	for _, item := range recvMsg.PropsList {
		totalAmount += item.TotalAmount
	}

	//实时更新日榜  直接通过db计算  1000个比邻币长一个亲密度, 所以gift_relation_val需要存float才行
	MysqlUpdateUserDailyRelationValue(recvMsg.ToUid, recvMsg.FromUid, 0, float32(totalAmount)/1000)

	//实时更新总榜
	{
		MysqlUpdateTotalRelationValue(recvMsg.ToUid, recvMsg.FromUid, 0, float32(totalAmount)/1000)
	}

	log.Info(prefix, zap.Any("recvMsg", recvMsg), zap.Any("totalAmount", totalAmount))
}

func ConsumerRabbitMQ() {
	const prefix = "ConsumerRabbitMQ "
	rabbitUrl := config.GetAppConfig().RabbitMqAddr
	if len(rabbitUrl) == 0 {
		panic("rabbitUrl not config!!!")
	}
	conn, err := amqp.Dial(rabbitUrl[0])
	if err != nil {
		log.Error(prefix, zap.Any("Failed to connect to RabbitMQ", err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Error(prefix, zap.Any("Failed to open a channel", err))
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Error(prefix, zap.Any("Failed to declare an exchange", err))
	}

	_, err = ch.QueueDeclare(
		QueueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Error(prefix, zap.Any("Failed to declare a queue", err))
	}

	err = ch.QueueBind(
		QueueName,    // queue name
		RoutingKey,   // routing key
		ExchangeName, // exchange
		false,
		nil)
	if err != nil {
		log.Error(prefix, zap.Any("Failed to bind a queue", err))
	}

	msgs, err := ch.Consume(
		QueueName, // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		log.Error(prefix, zap.Any("Failed to register a consumer", err))
	}

	log.Info(prefix + "begin")
	for d := range msgs {
		handleRabbitMqMsg(d.Body)
	}

	log.Error(prefix + "exit, why?")
}
