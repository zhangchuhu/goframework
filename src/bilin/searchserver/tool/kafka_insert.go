package main

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
)

func main() {
	kafkaBrokers := []string{
		"kafkawx001-core001.yy.com:8101",
		"kafkawx001-core002.yy.com:8101",
		"kafkawx001-core003.yy.com:8101",
	}
	kafkaTopic := os.Args[1]
	data := os.Args[2]

	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(kafkaBrokers, c)
	if err != nil {
		log.Fatalf("can not create kafka producer: %v", err)
	}

	_, _, err = p.SendMessage(&sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.ByteEncoder([]byte(data)),
	})
	if err != nil {
		log.Fatalf("can not send kafka message: %v", err)
	}
}
