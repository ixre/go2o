package msq

import (
	"fmt"
	"log"
	"time"

	"github.com/ixre/go2o/core/domain/interface/registry"
)

var (
	producer Producer
)

const NATS = 1
const KAFKA = 2
const REDIS = 3

type Producer interface {
	// Push 推送消息
	Push(topic string, message string) error
	// Close 关闭生产者
	Close()
}

// Configure 设置
func Configure(mqType int, address []string) error {
	if mqType == KAFKA {
		panic("if you want to use kafka as mq server. please uncomment blow line")
		//producer = newKafkaProducer(address)
	}
	if mqType == NATS {
		var err error
		producer, err = newNatsProducer("nats://" + address[0])
		return err
	}
	return fmt.Errorf("not implement mq type %d", mqType)
}

// 检查是否开启推送
func checkNatsSubs() bool {
	var repo registry.IRegistryRepo
	//repo := inject.GetRegistryRepo()
	v, _ := repo.GetValue(registry.AppEnableNatsSubscription)
	return v == "1" || v == "true"
}

// Push 推送消息
func Push(topic string, message string) error {
	if producer != nil {
		if checkNatsSubs() {
			return producer.Push(topic, message)
		}
	}
	log.Println("[ GO2O][ WARNING]: nats producer not available")
	return nil
}

// PushDelay 延迟推送消息
func PushDelay(topic string, message string, delay int) error {
	if producer == nil {
		return nil
	}
	if delay > 0 {
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
	return Push(topic, message)
}

// Close 关闭生产者
func Close() {
	if producer != nil {
		producer.Close()
	}
}
