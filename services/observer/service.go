package observer

import (
	"RabbitExample/services/counter/sdk"
	"context"
	"log"
	"time"
)

type Observer struct {
	counter *sdk.Counter
	cancel  context.CancelFunc
}

func New(c *sdk.Counter) *Observer {
	ctx, cancel := context.WithCancel(context.Background())
	o := Observer{
		counter: c,
		cancel:  cancel,
	}

	go o.run(ctx)

	return &o
}

func (o *Observer) Close() {
	o.cancel()
}

func (o *Observer) run(ctx context.Context) {
	log.Println("Observer service is running")
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second):
			log.Println("Count", count)
			ctx, cancel := context.WithTimeout(ctx, time.Second*15)
			c, err := o.counter.Increment(ctx, count)
			if err != nil {
				log.Println("Increment", err)
			}
			log.Println("Resp", c)
			cancel()
			count = c
		}
	}
}
