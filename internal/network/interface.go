package network

import "context"

type MessageHandlerFunc func(topic string, payload []byte)

type Broker interface {
	Connect() error
	Start(ctx context.Context) error
	Disconnect() error
	Subscribe(topic string, handler MessageHandlerFunc) error
	Publish(topic, addr string, payload []byte) error
	GetMessageChannel() <-chan Message
}

type Message struct {
	Topic   string
	Payload []byte
}
