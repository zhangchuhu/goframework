package updater

import (
	"bilin/searchserver/config"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/Shopify/sarama"
	kafka "github.com/bsm/sarama-cluster"
)

const (
	kafkaGroup = "searchserver-consumer-group"
)

var (
	httpClient *http.Client
)

func init() {
	httpDialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	httpTransport := &http.Transport{
		DialContext:       httpDialer.DialContext,
		DisableKeepAlives: false,
	}
	httpClient = &http.Client{
		Transport: httpTransport,
		Timeout:   5 * time.Second,
	}
}

type Updatable interface {
	ToURL(entity string) (*url.URL, error)
}

type UserU struct {
	Id      string `json:"id"`
	BilinId string `json:"bilin_id"`
	Name    string `json:"name"`
}

type RoomU struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Live      string `json:"live"`
	DisplayId string `json:"display_id"`
}

type SongU struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

func (u *UserU) ToURL(entity string) (*url.URL, error) {
	return convertToURL(entity, u)
}

func (r *RoomU) ToURL(entity string) (*url.URL, error) {
	return convertToURL(entity, r)
}

func (s *SongU) ToURL(entity string) (*url.URL, error) {
	return convertToURL(entity, s)
}

func convertToURL(entity string, obj Updatable) (l *url.URL, err error) {
	if l, err = url.ParseRequestURI(config.GetAppConfig().UpdateURL + "/" + entity); err != nil {
		return
	}
	l.RawQuery = convertToQuery(obj).Encode()
	return
}

func convertToQuery(obj Updatable) (query url.Values) {
	query = make(url.Values)
	val := reflect.ValueOf(obj).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		switch valueField.Kind() {
		case reflect.String:
			query.Add(typeField.Tag.Get("json"), valueField.String())
		default:
			panic("updatable object fields should be string")
		}
	}
	return
}

func KafkaLoop() {
	c := config.GetAppConfig()

	cc := kafka.NewConfig()
	cc.Consumer.Return.Errors = true
	cc.Group.Return.Notifications = true

	log.Info("new kafka consumer", zap.Any("topics", c.KafkaTopics), zap.Any("brokers", c.KafkaBrokers))
	consumer, err := kafka.NewConsumer(c.KafkaBrokers, kafkaGroup, c.KafkaTopics, cc)
	if err != nil {
		log.Fatal("can not create kafka consumer", zap.Error(err))
		panic(err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for err := range consumer.Errors() {
			log.Error("kafka consumer error", zap.Error(err))
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Warn("kafka consumer rebalanced", zap.Any("notice", ntf))
		}
	}()

	log.Info("consuming kafka messages, and watching signals")
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				log.Info("got kafka message",
					zap.String("topic", msg.Topic),
					zap.Int32("partition", msg.Partition),
					zap.Int64("offset", msg.Offset),
					zap.String("key", string(msg.Key)),
					zap.String("value", string(msg.Value)),
				)
				handleKafkaMessage(msg)
				consumer.MarkOffset(msg, "")
			}
		case <-signals:
			return
		}
	}
}

func handleKafkaMessage(msg *sarama.ConsumerMessage) {
	var (
		user UserU
		room RoomU
		song SongU
		link *url.URL
		err  error
		rsp  *http.Response
		body []byte
	)

	switch msg.Topic {
	case "bilin_user_update", "bilin_user_update_test":
		if err = json.Unmarshal(msg.Value, &user); err != nil {
			log.Error("can not parse", zap.Error(err), zap.String("user", string(msg.Value)))
			return
		}
		if link, err = user.ToURL(msg.Topic); err != nil {
			log.Error("can not convert to url", zap.Error(err), zap.Any("user", user))
			return
		}
	case "bilin_room_update", "bilin_room_update_test":
		if err = json.Unmarshal(msg.Value, &room); err != nil {
			log.Error("can not parse", zap.Error(err), zap.String("room", string(msg.Value)))
			return
		}
		if link, err = room.ToURL(msg.Topic); err != nil {
			log.Error("can not convert to url", zap.Error(err), zap.Any("room", room))
			return
		}
	case "bilin_song_update", "bilin_song_update_test":
		if err = json.Unmarshal(msg.Value, &song); err != nil {
			log.Error("can not parse", zap.Error(err), zap.String("song", string(msg.Value)))
			return
		}
		if link, err = song.ToURL(msg.Topic); err != nil {
			log.Error("can not convert to url", zap.Error(err), zap.Any("song", song))
			return
		}
	}

	if rsp, err = httpClient.Get(link.String()); err != nil {
		log.Error("can not do http get", zap.Error(err), zap.String("url", link.String()))
		return
	}
	defer rsp.Body.Close()
	if body, err = ioutil.ReadAll(rsp.Body); err != nil {
		log.Error("can not read http body", zap.Error(err), zap.String("url", link.String()))
		return
	}
	log.Info("update done",
		zap.String("topic", msg.Topic),
		zap.String("value", string(msg.Value)),
		zap.String("link", link.String()),
		zap.String("result", string(body)),
	)
}
