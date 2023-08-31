package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Connection *amqp.Connection
	Chanel     *amqp.Channel
	Queue      *amqp.Queue
	Deliveries <-chan amqp.Delivery
}

func NewRabbit(url, queue string, exclusive, autoAck bool) (*Rabbit, error) {
	r := Rabbit{}
	err := r.connect(url, queue, exclusive, autoAck)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Rabbit) connect(url, queueName string, exclusive, autoAck bool) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	r.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	/*if !exclusive {
		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
	}*/

	r.Chanel = ch

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		exclusive, // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	r.Queue = &q

	msg, err := r.Chanel.Consume(
		r.Queue.Name,
		"",
		autoAck,
		false,
		false,
		false,
		nil)

	if err != nil {
		return err
	}
	r.Deliveries = msg

	return nil
}

func (r *Rabbit) Close() {
	if r.Chanel != nil {
		r.Chanel.Close()
	}

	if r.Connection != nil {
		r.Connection.Close()
	}
}
