package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Dispatcher struct {
	broker    *Rabbit
	cancel    context.CancelFunc
	listeners map[string]chan []byte
	mx        sync.Mutex
	lstName   string
}

func NewDispatcher(url, lstName string) (*Dispatcher, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	r, err := NewRabbit(url, "", true, true)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := Dispatcher{
		broker:    r,
		cancel:    cancel,
		listeners: map[string]chan []byte{},
		lstName:   lstName,
	}

	go d.run(ctx)

	return &d, nil
}

func (d *Dispatcher) Close() {
	d.cancel()
	d.broker.Close()
}

func (d *Dispatcher) Dispatch(ctx context.Context, eventName string, data []byte) ([]byte, error) {
	corrId := randomString(32)
	ch := make(chan []byte)
	d.mx.Lock()
	d.listeners[corrId] = ch
	d.mx.Unlock()

	defer func() {
		close(ch)
		d.mx.Lock()
		delete(d.listeners, corrId)
		d.mx.Unlock()
	}()

	err := d.broker.Chanel.PublishWithContext(ctx,
		"",        // exchange
		d.lstName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Type:          eventName,
			CorrelationId: corrId,
			ReplyTo:       d.broker.Queue.Name,
			Body:          data,
		})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case response := <-ch:
		return response, nil
	}
}

func (d *Dispatcher) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-d.broker.Deliveries:
			d.mx.Lock()
			ch, ok := d.listeners[m.CorrelationId]
			d.mx.Unlock()

			if ok {
				ch <- m.Body
				err := m.Ack(false)
				if err != nil {
					log.Println("Ack error", err)
					break
				}
			}
		}
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
