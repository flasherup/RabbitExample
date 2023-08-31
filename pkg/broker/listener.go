package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type EventHandler func(request []byte) ([]byte, error)

type Listener struct {
	broker    *Rabbit
	cancel    context.CancelFunc
	listeners map[string]EventHandler
}

func NewListener(url, queue string) (*Listener, error) {
	r, err := NewRabbit(url, queue, false, false)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	l := Listener{
		broker:    r,
		cancel:    cancel,
		listeners: map[string]EventHandler{},
	}

	go l.run(ctx)

	return &l, nil
}

func (l *Listener) Close() {
	l.cancel()
	l.broker.Close()
}

func (l *Listener) Listen(eventName string, handler EventHandler) {
	l.listeners[eventName] = handler
}

func (l *Listener) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-l.broker.Deliveries:
			l.handleDelivery(&d)
		}
	}
}

func (l *Listener) handleDelivery(d *amqp.Delivery) {
	defer d.Ack(false)
	if handler, ok := l.listeners[d.Type]; ok {
		rsp, err := handler(d.Body)
		if err != nil {
			log.Println("Handler error", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = l.broker.Chanel.PublishWithContext(ctx,
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				Type:          d.Type,
				CorrelationId: d.CorrelationId,
				Body:          rsp,
			})
		if err != nil {
			log.Println("Reply error", err)
		}
	}
}
