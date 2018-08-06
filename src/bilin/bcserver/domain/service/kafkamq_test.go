package service

import (
	"testing"
	"bilin/bcserver/config"
)

func TestKafkaProducerInit(t *testing.T) {
	config.SetTestAppConfig(mysqlconfig)
	KafkaProducerInit()
}
