package main

import (
	"RabbitExample/pkg/broker"
	"RabbitExample/services/counter"
	"RabbitExample/services/counter/sdk"
)

func main() {
	l, err := broker.NewListener("amqp://admin:admin@localhost:5672/", sdk.QueueName)
	if err != nil {
		panic(err)
	}

	var forever chan struct{}
	s := counter.New(l)
	defer s.Close()
	<-forever
}
