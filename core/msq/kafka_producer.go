package msq

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

var _ Producer = new(KafkaProducer)

type KafkaProducer struct {
	pro     sarama.AsyncProducer
	address []string
}

func newKafkaProducer(address []string) *KafkaProducer {
	k := &KafkaProducer{
		address: address,
		pro:     createKafkaProducer(address),
	}
	return k
}

// 创建异步producer
func createKafkaProducer(address []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoResponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestamp没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_11_0_2
	config.ClientID = "go2o-server-producer"
	log.Println("[ Go2o][ Info]: start kafka producer")
	//使用配置,新建一个异步生产者
	var producer sarama.AsyncProducer
	var err error
	var retryTimes = 10
	log.Println("[ Go2o][ Kafka]: start kafka producer...")
	for i:=0 ;i<retryTimes;i++ {
		producer, err = sarama.NewAsyncProducer(address, config)
		if err == nil {
			break
		} else if i == retryTimes-1 {
			log.Println("[ Go2o][ Kafka]: can't connect kafka server! ", err.Error())
			os.Exit(1)
		}
		time.Sleep(time.Second * 2)
	}

	//defer producer.AsyncClose()
	// 判断哪个通道发送过来数据.
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case suc := <-p.Successes():
				if suc != nil {
					//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				}
			case fail := <-p.Errors():
				if fail != nil {
					fmt.Println("err: ", fail.Err)
				}
			}
		}
	}(producer)
	return producer
}

func (k *KafkaProducer) Push(topic string, key string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	if len(key) > 0 {
		msg.Key = sarama.ByteEncoder(key)
	}
	//使用通道发送
	k.pro.Input() <- msg
	return nil
}

func (k *KafkaProducer) Close() {
	k.pro.AsyncClose()
}
