package msq

import (
	nats "github.com/nats-io/nats.go"
	"log"
)

var _ Producer = new(natsProducer)

type natsProducer struct {
	nc      *nats.Conn
	address string
}

func newNatsProducer(address string) (Producer, error) {
	log.Println("[ Go2o][ Mq]: start nats producer...")
	nc, err := nats.Connect(address)
	if err != nil {
		log.Println("[ Go2o][ Mq]: can't connect nats server!", err.Error())
		return nil, err
	}
	return &natsProducer{
		address: address,
		nc:      nc,
	}, nil
}

func (n natsProducer) Push(topic string, message string) error {
	return n.nc.Publish(topic, []byte(message))
}

func (n natsProducer) Close() {
	n.nc.Close()
}
