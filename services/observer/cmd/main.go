package main

import (
	"RabbitExample/pkg/broker"
	"RabbitExample/services/counter/sdk"
	"RabbitExample/services/observer"
)

func main() {
	d, err := broker.NewDispatcher("amqp://admin:admin@localhost:5672/", sdk.QueueName)
	if err != nil {
		panic(err)
	}

	cs := sdk.New(d)

	var forever chan struct{}
	s := observer.New(cs)
	defer s.Close()
	<-forever
}
