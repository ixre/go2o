package msq

import (
	"errors"
	"fmt"
	"time"
)

var (
	producer Producer
)

const NATS = 1
const KAFKA = 2
const REDIS = 3

type Producer interface {
	// 推送消息
	Push(topic string, message string) error
	// 关闭生产者
	Close()
}

// 设置
func Configure(mqType int, address []string) error {
	if mqType == KAFKA {
		panic("if you want to use kafka as mq server. please uncomment blow line")
		//producer = newKafkaProducer(address)
		return nil
	}
	if mqType == NATS {
		var err error
		producer, err = newNatsProducer("nats://" + address[0])
		return err
	}
	return errors.New(fmt.Sprintf("not implement mq type %d", mqType))
}

// 推送消息
func Push(topic string, message string) error {
	if producer != nil {
		return producer.Push(topic, message)
	}
	return nil
}

// 延迟推送消息
func PushDelay(topic string, message string, delay int) error {
	if producer == nil {
		return nil
	}
	if delay > 0 {
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
	return Push(topic, message)
}

// 关闭生产者
func Close() {
	if producer != nil {
		producer.Close()
	}
}
