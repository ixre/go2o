package msq

import (
	"github.com/ixre/gof/log"
	"github.com/nats-io/nats.go"
)

var _ Producer = new(natsProducer)

type natsProducer struct {
	nc      *nats.Conn
	address string
}

func newNatsProducer(address string) (Producer, error) {
	log.Println("[ GO2O][ INFO]: start nats producer...")
	nc, err := nats.Connect(address)
	if err != nil {
		log.Println("[ GO2O][ ERROR]: can't connect nats server!", err.Error())
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
