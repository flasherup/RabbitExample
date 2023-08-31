package sdk

import (
	"RabbitExample/services/counter/model"
	"context"
	"encoding/json"
)

const (
	QueueName = "rpc_queue"
)

type dispatcher interface {
	Dispatch(ctx context.Context, eventName string, request []byte) ([]byte, error)
	Close()
}

type Counter struct {
	eventDispatcher dispatcher
}

func New(d dispatcher) *Counter {
	c := Counter{
		eventDispatcher: d,
	}
	return &c
}

func (c *Counter) Close() {
	c.eventDispatcher.Close()
}

func (c *Counter) Increment(ctx context.Context, count int) (int, error) {
	req := model.CountReq{
		Count: count,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	rsp, err := c.eventDispatcher.Dispatch(ctx, "increment", b)
	if err != nil {
		return 0, err
	}

	r := model.CountRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return 0, err
	}

	return r.Count, nil
}
